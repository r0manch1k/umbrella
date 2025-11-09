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

#include <openssl/pem.h>
#include <openssl/err.h>
#include <openssl/rsa.h>
#include <openssl/sha.h>

LicenseManager::LicenseManager(QObject *parent)
    : QObject(parent)
{
    qDebug() << "LicenseManager constructor called";
    networkManager = new QNetworkAccessManager(this);
}

void LicenseManager::requestLicenseFromServer(const QString &userId, int duration, const QString &hwFingerprint = QString())
{
    qDebug() << "requestLicenseFromServer called with userId:" << userId << "duration:" << duration << "hwFingerprint:" << hwFingerprint;
    QUrl url("http://127.0.0.1:5000/license/issue");
    QNetworkRequest request(url);
    request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject obj;
    obj["user_id"] = userId;
    obj["duration"] = duration;
    if (!hwFingerprint.isEmpty())
        obj["hw_fingerprint"] = hwFingerprint;

    QJsonDocument doc(obj);
    QByteArray data = doc.toJson(QJsonDocument::Compact);

    qDebug() << "Requesting license from server:";
    qDebug() << "URL:" << url.toString();
    qDebug() << "Request JSON:" << QString::fromUtf8(data);

    QNetworkReply *reply = networkManager->post(request, data);

    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        qDebug() << "License request finished.";
        qDebug() << "HTTP status code:" << reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
        qDebug() << "Network reply error:" << reply->error() << "-" << reply->errorString();

        QByteArray responseData = reply->readAll();
        qDebug() << "Response body:" << QString::fromUtf8(responseData);

        QJsonDocument respDoc = QJsonDocument::fromJson(responseData);
        if (respDoc.isNull()) {
            qDebug() << "Response is not a valid JSON document.";
        } else if (!respDoc.isObject()) {
            qDebug() << "Response JSON is not an object.";
        } else {
            QJsonObject respObj = respDoc.object();
            qDebug() << "Response JSON object keys:" << respObj.keys();
            if (!respObj.contains("license")) {
                qDebug() << "Response JSON does not contain 'license' key.";
            }
            if (!respObj.contains("signature")) {
                qDebug() << "Response JSON does not contain 'signature' key.";
            }
            if (reply->error() == QNetworkReply::NoError) {
                QString licenseJson = respObj["license"].toString();
                QByteArray signature = QByteArray::fromBase64(respObj["signature"].toString().toUtf8());
                qDebug() << "License JSON received:" << licenseJson;
                qDebug() << "Signature (base64 decoded):" << signature.toHex();

                activateFromServerResponse(licenseJson, signature);
            }
        }

        reply->deleteLater();
    });
}

void LicenseManager::verifyLicenseWithServer()
{
    qDebug() << "verifyLicenseWithServer called";
    if (m_licenseData.isEmpty() || m_signature.isEmpty()) {
        qDebug() << "No license data or signature to verify";
        return;
    }

    QUrl url("http://127.0.0.1:5000/license/verify");
    QNetworkRequest request(url);
    request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject obj;
    obj["license"] = QString::fromUtf8(m_licenseData);
    obj["signature"] = QString::fromUtf8(m_signature.toBase64());

    qDebug() << "Sending verification request with license:" << obj["license"];
    qDebug() << "Sending verification request with signature (base64):" << obj["signature"];

    QNetworkReply *reply = networkManager->post(request, QJsonDocument(obj).toJson());

    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QByteArray responseData = reply->readAll();
            qDebug() << "Verification response data:" << QString::fromUtf8(responseData);
            QJsonDocument doc = QJsonDocument::fromJson(responseData);
            if (!doc.isNull() && doc.object().contains("valid")) {
                bool valid = doc.object()["valid"].toBool();
                m_valid = valid;
                qDebug() << "Server license check result:" << valid;
            } else {
                qDebug() << "Verification response JSON invalid or missing 'valid' key";
            }
        } else {
            qDebug() << "License verify error:" << reply->errorString();
        }
        reply->deleteLater();
    });
}

QByteArray LicenseManager::loadPublicKey() const
{
    qDebug() << "loadPublicKey called";
    QFile f(":/keys/public.pem");
    if (!f.open(QIODevice::ReadOnly))
    {
        qWarning() << "Cannot open public.pem";
        return {};
    }
    QByteArray keyData = f.readAll();
    qDebug() << "Public key loaded, size:" << keyData.size();
    return keyData;
}

bool LicenseManager::verifySignature(const QByteArray &payload, const QByteArray &sig)
{
    qDebug() << "verifySignature called with payload size:" << payload.size() << "signature size:" << sig.size();
    QByteArray pubKeyData = loadPublicKey();
    if (pubKeyData.isEmpty()) {
        qDebug() << "Public key data is empty, cannot verify signature";
        return false;
    }

    BIO *bio = BIO_new_mem_buf(pubKeyData.constData(), pubKeyData.size());
    RSA *rsa = PEM_read_bio_RSA_PUBKEY(bio, nullptr, nullptr, nullptr);
    BIO_free(bio);

    if (!rsa) {
        qDebug() << "Failed to load RSA public key";
        return false;
    }

    bool result = false;
    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256(reinterpret_cast<const unsigned char*>(payload.constData()), payload.size(), hash);
    qDebug() << "SHA256 hash computed for payload";

    int rc = RSA_verify(NID_sha256, hash, SHA256_DIGEST_LENGTH,
                        reinterpret_cast<const unsigned char*>(sig.constData()),
                        sig.size(), rsa);
    if (rc == 1) {
        result = true;
        qDebug() << "Signature verification succeeded";
    } else {
        qDebug() << "Signature verification failed";
    }

    RSA_free(rsa);
    return result;
}

void LicenseManager::saveLicense(const QString &licenseJson)
{
    qDebug() << "saveLicense called with licenseJson size:" << licenseJson.size();
    QFile f("license.json");
    if (f.open(QIODevice::WriteOnly))
    {
        f.write(licenseJson.toUtf8());
        f.close();
        qDebug() << "License saved to file license.json";
    }
    else {
        qDebug() << "Failed to open license.json for writing";
    }

    m_licenseData = licenseJson.toUtf8();
    QJsonDocument doc = QJsonDocument::fromJson(m_licenseData);
    if (!doc.isNull())
    {
        QJsonObject obj = doc.object();
        if (obj.contains("expires")) {
            m_expiration = QDateTime::fromString(obj["expires"].toString(), Qt::ISODate);
            qDebug() << "License expiration set to:" << m_expiration.toString(Qt::ISODate);
        } else {
            qDebug() << "License JSON does not contain 'expires' field";
        }
    }
    else {
        qDebug() << "License JSON is invalid";
    }
    m_valid = true;
    qDebug() << "License marked as valid";
}

bool LicenseManager::isLicenseValid() const
{
    qDebug() << "isLicenseValid called, m_valid:" << m_valid << "isExpired():" << isExpired();
    return m_valid && !isExpired();
}

bool LicenseManager::isExpired() const
{
    qDebug() << "isExpired called, expiration date:" << m_expiration.toString(Qt::ISODate);
    if (!m_expiration.isValid()) {
        qDebug() << "Expiration date is invalid, license considered expired";
        return true;
    }
    bool expired = QDateTime::currentDateTimeUtc() > m_expiration;
    qDebug() << "Current UTC datetime:" << QDateTime::currentDateTimeUtc().toString(Qt::ISODate) << "expired:" << expired;
    return expired;
}

void LicenseManager::activateFromServerResponse(const QString &licenseJson, const QByteArray &signature)
{
    qDebug() << "activateFromServerResponse called";
    qDebug() << "License JSON size:" << licenseJson.size();
    qDebug() << "Signature size:" << signature.size();

    if (verifySignature(licenseJson.toUtf8(), signature))
    {
        qDebug() << "Signature verified successfully, saving license";
        saveLicense(licenseJson);
        m_signature = signature;
        m_valid = true;
        qDebug() << "License activated and marked valid";
    }
    else
    {
        m_valid = false;
        qDebug() << "Signature verification failed, license not activated";
    }
}

QString LicenseManager::currentLicense() const
{
    qDebug() << "currentLicense called, license data size:" << m_licenseData.size();
    return QString::fromUtf8(m_licenseData);
}