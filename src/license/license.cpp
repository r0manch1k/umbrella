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

void LicenseManager::verify()
{
    qDebug() << "verify";

    if (m_l.isEmpty() || m_s.isEmpty())
    {
        qDebug() << "license|signature is empty";
    }

    QNetworkRequest req;
    req.setUrl(QUrl("http://127.0.0.1:9090/license/verify"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "text/plain");

    QJsonObject reqjo;
    reqjo["license"] = QString::fromUtf8(m_l);
    reqjo["signature"] = QString::fromUtf8(m_s);

    QByteArray j = QJsonDocument(reqjo).toJson(QJsonDocument::Compact);

    QByteArray enc = this->enc(j);

    qDebug() << "req license:" << reqjo["license"];
    qDebug() << "req signature:" << reqjo["signature"];

    QNetworkReply *res = nm->post(req, enc);

    connect(res, &QNetworkReply::finished, this, [this, res]()
            {
        if (res->error() == QNetworkReply::NoError)
        {
            QByteArray resba = res->readAll();
            qDebug() << "response data:" << QString::fromUtf8(resba);
            QByteArray resdec;
            try
            {
                resdec = this->dec(resba);
            }
            catch (...)
            {
                qDebug() << "reply is not valid Base64!";
                return;
            }
            qDebug() << "decoded" << resdec.toStdString();
            QJsonParseError e;
            QJsonDocument resjd = QJsonDocument::fromJson(resdec, &e);
            if (e.error != QJsonParseError::NoError)
            {
                qDebug() << "❌ JSON parse error:" << e.errorString();
                return;
            }

            if (!resjd.isObject())
            {
                qDebug() << "❌ Reply is not a JSON object";
                return;
            }
            if (!resjd.isNull() && resjd.object().contains("payload") && resjd.object().contains("signature"))
            {
                QJsonObject resjo = resjd.object();
                if (!resjo.contains("payload") || !resjo.contains("signature"))
                {
                    qDebug() << "Missing payload or payload_signature";
                    m_v = false;
                    return;
                }
                if (resjo["payload"].isBool() && !resjo["payload"].toBool())
                {
                    qDebug() << "payload is bool";
                    m_v = false;
                }
                else if (resjo["payload"].isObject())
                {
                    QJsonDocument p(resjo["payload"].toObject());
                    QByteArray pba = p.toJson(QJsonDocument::Compact);
                    QJsonDocument pjd = QJsonDocument::fromJson(pba);
                    if (pjd.isObject() && pjd.object().contains("signature"))
                    {
                        QString sig = pjd.object()["signature"].toString();
                        if (issigned(QByteArray::fromBase64(pjd.object()["license_value"].toString().toUtf8()), sig.toUtf8()))
                        {
                            m_s = sig.toUtf8();
                            save(m_l, m_s);
                            load();
                            m_v = true;
                            qDebug() << "License verified and saved successfully";
                            return;
                        }
                        else
                        {
                            m_v = false;
                        }
                        res->deleteLater();
                        qDebug() << "license check result:" << m_v;
                    }
                    else
                    {
                        qDebug() << "no valid key";
                    }
                }
                else
                {
                    qDebug() << "error:" << res->errorString();
                }
                m_v = false;
                res->deleteLater();
            } else {
                qDebug() << "payload error:" << res->error();
            }
        } else {
            qDebug() << "res error:" << res->error();
        } });
}

void LicenseManager::save(const QByteArray &l, const QByteArray &s)
{
    qDebug() << "save";
    qDebug() << "license size:" << l.size() << "signature size:" << s.size();
    QFile f("_");
    if (f.open(QIODevice::WriteOnly))
    {
        QDataStream out(&f);
        out << quint32(s.size());
        if (!s.isEmpty())
            out.writeRawData(s.constData(), s.size());
        out << quint32(l.size());
        if (!l.isEmpty())
            out.writeRawData(s.constData(), s.size());
        f.close();
        qDebug() << "license saved";
    }
    else
    {
        qDebug() << "license saved";
    }

    m_l = l;
    if (!s.isEmpty())
        m_s = s;

    qDebug() << "saved";
}

void LicenseManager::load()
{
    qDebug() << "load";
    QFile f("_");
    if (!f.exists())
    {
        qDebug() << ". not found";
        return;
    }
    if (f.open(QIODevice::ReadOnly))
    {
        QDataStream in(&f);
        quint32 siglen = 0;
        in >> siglen;
        QByteArray s;
        if (siglen > 0)
        {
            s.resize(siglen);
            in.readRawData(s.data(), siglen);
        }
        quint32 liclen = 0;
        in >> liclen;
        QByteArray l;
        if (liclen > 0)
        {
            l.resize(liclen);
            in.readRawData(l.data(), liclen);
        }
        f.close();

        qDebug() << "Loaded license.bin, license size:" << l.size() << "signature size:" << s.size();

        if (!l.isEmpty())
        {
            if (!s.isEmpty())
            {
                if (issigned(l, s))
                {
                    m_l = l;
                    m_s = s;
                    qDebug() << "Loaded license: signature valid";
                    return;
                }
            }
            qDebug() << "Loaded license: signature invalid";
        }
    }
    else
    {
        qDebug() << "Failed to open license.bin for reading";
    }
}

bool LicenseManager::issigned(const QByteArray &l, const QByteArray &s)
{
    qDebug() << "issigned() — verifying license using embedded public.pem" << l.toStdString() << s.toStdString();

    QFile f(":/keys/public.pem");
    if (!f.open(QIODevice::ReadOnly))
    {
        qDebug() << "Cannot open embedded public.pem";
        return false;
    }
    QByteArray keyData = f.readAll();
    f.close();

    BIO *bio = BIO_new_mem_buf(keyData.data(), keyData.size());
    if (!bio)
    {
        qDebug() << "BIO_new_mem_buf failed";
        return false;
    }

    RSA *rsa = PEM_read_bio_RSA_PUBKEY(bio, nullptr, nullptr, nullptr);
    BIO_free(bio);

    if (!rsa)
    {
        qDebug() << "PEM_read_bio_RSA_PUBKEY failed:" << ERR_error_string(ERR_get_error(), nullptr);
        return false;
    }

    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256(reinterpret_cast<const unsigned char *>(l.constData()), l.size(), hash);

    QByteArray sba = QByteArray::fromBase64(s);

    int res = RSA_verify(NID_sha256,
                         hash, SHA256_DIGEST_LENGTH,
                         reinterpret_cast<const unsigned char *>(sba.constData()), sba.size(),
                         rsa);
    if (res <= 0)
        qDebug() << "RSA_verify failed:" << ERR_error_string(ERR_get_error(), nullptr);

    RSA_free(rsa);

    return res > 0;
}

void LicenseManager::license(const QByteArray &l)
{
    m_l = l;
    qDebug() << "currentLicense called, license data size:" << m_l.size();
    return;
}

QByteArray LicenseManager::enc(const QByteArray &l) const
{
    return l.toBase64();
}

QByteArray LicenseManager::dec(const QByteArray &l) const
{
    return QByteArray::fromBase64(l);
}