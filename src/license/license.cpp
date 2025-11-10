#include "license.h"

#include <QFile>
#include <QJsonDocument>
#include <QJsonObject>
#include <QDebug>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QJsonObject>
#include <QJsonDocument>

#include <QDataStream>

#include <openssl/pem.h>
#include <openssl/rsa.h>
#include <openssl/sha.h>
#include <openssl/err.h>


LicenseManager::LicenseManager(QObject *parent)
    : QObject(parent)
{
    nm = new QNetworkAccessManager(this);
    load();
}

void LicenseManager::issue(const QString &userId, int duration_hours, const QString &hwFingerprint = QString())
{
    qDebug() << "issue";

    QNetworkRequest req;
    req.setUrl(QUrl("http://127.0.0.1:5000/license/issue"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject reqo;
    reqo["user_id"] = userId;
    reqo["duration_hours"] = duration_hours;
    if (!hwFingerprint.isEmpty())
        reqo["hw_fingerprint"] = hwFingerprint;

    QJsonDocument reqd(reqo);
    QByteArray reqba = reqd.toJson(QJsonDocument::Compact);

    qDebug() << "post /license/isuue";
    qDebug() << "json:" << QString::fromUtf8(reqba);

    QNetworkReply *res = nm->post(req, reqba);

    connect(res, &QNetworkReply::finished, this, [this, res]() {
        qDebug() << "response";
        qDebug() << "status:" << res->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
        qDebug() << "error:" << res->error() << "-" << res->errorString();

        QByteArray resba = res->readAll();
        qDebug() << "body:" << QString::fromUtf8(resba);

        QJsonDocument resd = QJsonDocument::fromJson(resba);
        if (resd.isNull()) {
            qDebug() << "null";
        } else if (!resd.isObject()) {
            qDebug() << "not an object.";
        } else {
            QJsonObject reso = resd.object();
            qDebug() << "response json fields" << reso.keys();
            if (!reso.contains("license")) {
                qDebug() << "no license key";
            }
            if (!reso.contains("signature")) {
                qDebug() << "no signature key";
            }
            if (res->error() == QNetworkReply::NoError) {
                QString l = reso["license"].toString();
                QByteArray s = QByteArray::fromBase64(reso["signature"].toString().toUtf8());
                qDebug() << "license" << l;
                qDebug() << "signature" << s.toHex();

                // activateFromServerResponse(l, s);
            }
        }

        res->deleteLater();
    });
}

void LicenseManager::verify()
{
    qDebug() << "verify";

    if (m_l.isEmpty() || m_s.isEmpty()) {
        qDebug() << "no license";
        return;
    }

    QNetworkRequest req;
    req.setUrl(QUrl("http://127.0.0.1:5000/license/verify"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject reqo;
    reqo["license"] = QString::fromUtf8(m_l);
    reqo["signature"] = QString::fromUtf8(m_s.toBase64());

    qDebug() << "req license:" << reqo["license"];
    qDebug() << "req signature (base64):" << reqo["signature"];

    QNetworkReply *res = nm->post(req, QJsonDocument(reqo).toJson());

    connect(res, &QNetworkReply::finished, this, [this, res]() {
        if (res->error() == QNetworkReply::NoError) {
            QByteArray resba = res->readAll();
            qDebug() << "response data:" << QString::fromUtf8(resba);
            QJsonDocument resd = QJsonDocument::fromJson(resba);
            if (!resd.isNull() && resd.object().contains("valid")) {
                bool v = resd.object()["valid"].toBool();
                // m_v = v;
                qDebug() << "license check result:" << v;
            } else {
                qDebug() << "no valid key";
            }
        } else {
            qDebug() << "error:" << res->errorString();
        }
        res->deleteLater();
    });
}


void LicenseManager::save(const QByteArray &l, const QByteArray &s)
{
    qDebug() << "save";
    qDebug() << "license size:" << l.size() << "signature size:" << s.size();
    QFile f(".l");
    if (f.open(QIODevice::WriteOnly))
    {
        QDataStream out(&f);
        out << quint32(s.size());
        if (!s.isEmpty()) out.writeRawData(s.constData(), s.size());
        out << quint32(l.size());
        if (!l.isEmpty()) out.writeRawData(s.constData(), s.size());
        f.close();
        qDebug() << "saved to .l";
    }
    else {
        qDebug() << "saved to .l";
    }

    m_l = l;
    if (!s.isEmpty()) m_s = s;

    qDebug() << "saved";
}

void LicenseManager::load()
{
    qDebug() << "load";
    QFile f(".l");
    if (!f.exists()) {
        qDebug() << ".l not found";
        return;
    }
    if (f.open(QIODevice::ReadOnly))
    {
        QDataStream in(&f);
        quint32 siglen = 0;
        in >> siglen;
        QByteArray s;
        if (siglen > 0) {
            s.resize(siglen);
            in.readRawData(s.data(), siglen);
        }
        quint32 liclen = 0;
        in >> liclen;
        QByteArray l;
        if (liclen > 0) {
            l.resize(liclen);
            in.readRawData(l.data(), liclen);
        }
        f.close();

        qDebug() << "Loaded license.bin, license size:" << l.size() << "signature size:" << s.size();

        if (!l.isEmpty()) {
            m_l = l;
            if (!s.isEmpty()) m_s = s;

            // Optionally, verify signature locally if verifySignature is available
            if (!m_s.isEmpty()) {
                if (issigned(m_l, m_s)) {
                    // m_valid = true;
                    qDebug() << "Loaded license: signature valid";
                } else {
                    // m_valid = false;
                    qDebug() << "Loaded license: signature invalid";
                }
            }
        }
    }
    else {
        qDebug() << "Failed to open license.bin for reading";
    }
}

bool LicenseManager::issigned(const QByteArray &l, const QByteArray &s)
{
    qDebug() << "issigned() — verifying license using embedded public.pem";

    QFile f(":/keys/public.pem");
    if (!f.open(QIODevice::ReadOnly)) {
        qDebug() << "Cannot open embedded public.pem";
        return false;
    }
    QByteArray keyData = f.readAll();
    f.close();

    BIO *bio = BIO_new_mem_buf(keyData.data(), keyData.size());
    if (!bio) {
        qDebug() << "BIO_new_mem_buf failed";
        return false;
    }

    RSA *rsa = PEM_read_bio_RSA_PUBKEY(bio, nullptr, nullptr, nullptr);
    BIO_free(bio);

    if (!rsa) {
        qDebug() << "PEM_read_bio_RSA_PUBKEY failed:" << ERR_error_string(ERR_get_error(), nullptr);
        return false;
    }

    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256(reinterpret_cast<const unsigned char*>(l.constData()), l.size(), hash);

    int res = RSA_verify(NID_sha256,
                         hash, SHA256_DIGEST_LENGTH,
                         reinterpret_cast<const unsigned char*>(s.constData()), s.size(),
                         rsa);

    RSA_free(rsa);

    return res > 0;
}


QString LicenseManager::license() const
{
    qDebug() << "currentLicense called, license data size:" << m_l.size();
    return QString::fromUtf8(m_l);
}