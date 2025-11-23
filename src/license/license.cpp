#include "license.h"

#include <openssl/err.h>
#include <openssl/pem.h>
#include <openssl/rsa.h>
#include <openssl/sha.h>

#include <QCoreApplication>
#include <QCryptographicHash>
#include <QDataStream>
#include <QDebug>
#include <QEventLoop>
#include <QFile>
#include <QHostInfo>
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkAccessManager>
#include <QNetworkInterface>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QScopedPointer>
#include <QScopedPointerDeleteLater>
#include <QSysInfo>
#include <vector>

LicenseManager::LicenseManager(QObject* parent) : QObject(parent) {
    std::vector<int> base64_public_key = {
        1309 ^ 123, 1077 ^ 123, 1089 ^ 123, 1086 ^ 123, 1405 ^ 123, 8572 ^ 123, 7691 ^ 123,
        291 ^ 123,  1086 ^ 123, 1088 ^ 123, 1086 ^ 123, 1088 ^ 123, 1072 ^ 123, 8572 ^ 123,
        8572 ^ 123, 116 ^ 123,  7715 ^ 123, 1077 ^ 123, 1089 ^ 123, 7715 ^ 123, 1072 ^ 123,
        8575 ^ 123, 1088 ^ 123, 1072 ^ 123, 501 ^ 123,  1400 ^ 123, 1077 ^ 123, 1110 ^ 123,
        1400 ^ 123, 1057 ^ 123, 1072 ^ 123, 8572 ^ 123, 1110 ^ 123, 102 ^ 123,  1086 ^ 123,
        114 ^ 123,  1400 ^ 123, 1110 ^ 123, 1072 ^ 123};

    for (int i = 0; i < base64_public_key.size(); ++i) {
        int j = base64_public_key[i] ^ 123;
        j = (j + i * 17) ^ (i * 31);
        j = (j ^ (i * 31)) - i * 17;
        b64pk.append(QChar(j));
    }

    net = new QNetworkAccessManager(this);
}

void LicenseManager::feed(uint32_t v, int n) {
    qDebug() << "void LicenseManager::feed" << v << n;
    v = (v << ((n % 3) + 1)) | (v >> (32 - ((n % 3) + 1)));
    h ^= (v + 0x9E3779B9u + n);
    h = (h * 0x85EBCA6Bu) ^ (h >> 13);
}

uint32_t LicenseManager::geth() {
    qDebug() << "H" << h;
    feed(b64pk[2].unicode(), 2);
    return h;
}

static inline uint32_t rotl(uint32_t x, int n) { return (x << n) | (x >> (32 - n)); }
static inline uint32_t rotr(uint32_t x, int n) { return (x >> n) | (x << (32 - n)); }

void LicenseManager::hash() {
    uint32_t orig = h;
    uint32_t a = orig ^ 0xDEADBEEFu;
    uint32_t b = a + 0x13579BDFu;
    uint32_t c = rotl(b, 7);
    uint32_t d = c ^ 0xCAFEBABEu;
    uint32_t e = d + 0xABCDEF01u;
    uint32_t f = rotl(e, 11);
    uint32_t g = f ^ 0xA5A5A5A5u;
    uint32_t htmp = rotl(g, 5);
    uint32_t i = htmp + 0x11223344u;
    uint32_t j = i ^ 0x55AA55AAu;
    uint32_t k = rotl(j, 3);
    uint32_t k2 = rotr(k, 3);
    uint32_t j2 = k2 ^ 0x55AA55AAu;
    uint32_t i2 = j2 - 0x11223344u;
    uint32_t htmp2 = rotr(i2, 5);
    uint32_t g2 = htmp2 ^ 0xA5A5A5A5u;
    uint32_t f2 = rotr(g2, 11);
    uint32_t e2 = f2 - 0xABCDEF01u;
    uint32_t d2 = e2 ^ 0xCAFEBABEu;
    uint32_t c2 = rotr(d2, 7);
    uint32_t b2 = c2 - 0x13579BDFu;
    uint32_t a2 = b2 ^ 0xDEADBEEFu;
    h = a2;
}

QString LicenseManager::fingerprint() {
    QString data;
    data += QSysInfo::machineHostName();
    data += QSysInfo::kernelType();
    data += QSysInfo::kernelVersion();
    data += QSysInfo::prettyProductName();
    for (const QNetworkInterface& iface : QNetworkInterface::allInterfaces()) {
        if (!(iface.flags() & QNetworkInterface::IsLoopBack) &&
            iface.flags() & QNetworkInterface::IsUp) {
            data += iface.hardwareAddress();
            break;
        }
    }
    QByteArray hash = QCryptographicHash::hash(data.toUtf8(), QCryptographicHash::Sha256);
    return hash.left(8).toHex();
}

void LicenseManager::verify() {
    qDebug() << "void LicenseManager::verify";

    if (ml.isEmpty()) {
        qDebug() << ml.isEmpty();
        return;
    }

    QNetworkRequest req;

    // req.setUrl(QUrl("https://thesokolovsky.ru/umbrella/license/verify"));
    req.setUrl(QUrl("http://localhost:9090/v1/license/verify"));

    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject req_json_obj;
    req_json_obj["license"] = QString::fromUtf8(ml);
    req_json_obj["fingerprint"] = fingerprint();

    QByteArray payload = QJsonDocument(req_json_obj).toJson(QJsonDocument::Compact);

    qDebug() << "req_json_obj[license]" << req_json_obj["license"];
    qDebug() << "req_json_obj[fingerprint]" << req_json_obj["fingerprint"];

    QNetworkReply* res;
    QJsonObject req_json_obj_fin;
    req_json_obj_fin["secret_payload"] = QString::fromUtf8(QEncodeByteArray(payload));
    QByteArray data = QJsonDocument(req_json_obj_fin).toJson(QJsonDocument::Compact);

    qDebug() << "req_json_obj_fin[secret_payload]" << req_json_obj_fin["secret_payload"];

    res = net->post(req, data);

    QEventLoop loop;
    QObject::connect(res, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    QScopedPointer<QNetworkReply, QScopedPointerDeleteLater> guard(res);

    if (res->error() != QNetworkReply::NoError) {
        if (res->error() == QNetworkReply::NetworkSessionFailedError ||
            res->error() == QNetworkReply::HostNotFoundError ||
            res->error() == QNetworkReply::TimeoutError ||
            res->error() == QNetworkReply::ConnectionRefusedError ||
            res->error() == QNetworkReply::UnknownNetworkError) {
            if (!ms.isEmpty() && issigned(ml, ms)) {
                save(ml, ms);
            }
            return;
        }
        qDebug() << "res->error() != QNetworkReply::NoError" << res->error();
        return;
    }
    QByteArray res_bytes = res->readAll();

    qDebug() << "QString::fromUtf8(res_bytes)" << QString::fromUtf8(res_bytes);

    QByteArray dec;
    try {
        dec = this->QDecodeByteArray(res_bytes);
    } catch (...) {
        qDebug() << "reply is not valid Base64!";
        return;
    }

    qDebug() << "dec.toStdString()" << dec.toStdString();

    QJsonParseError e;
    QJsonDocument res_json_doc = QJsonDocument::fromJson(dec, &e);
    if (e.error != QJsonParseError::NoError) {
        qDebug() << "e.error != QJsonParseError::NoError" << e.errorString();
        return;
    }

    if (!res_json_doc.isObject()) {
        qDebug() << "!res_json_doc.isObject()";
        return;
    }

    if (res_json_doc.isNull()) {
        qDebug() << "res_json_doc.isNull()";
        return;
    }

    QJsonObject res_json_obj = res_json_doc.object();

    if (!res_json_obj.contains("valid") || !res_json_obj.contains("signature")) {
        qDebug() << "!res_json_obj.contains(valid) || !res_json_obj.contains(signature)";
        return;
    }

    if (res_json_obj["valid"].isBool() && !res_json_obj["valid"].toBool()) {
        qDebug() << "res_json_obj[valid].isBool() && !res_json_obj[valid].toBool()";
        return;
    }

    if (!issigned(ml, res_json_obj["signature"].toString().toUtf8())) {
        qDebug() << "!issigned(ml, res_json_obj[signature].toString().toUtf8())";
        qDebug() << "res_json_obj[signature].toString().toUtf8()"
                 << res_json_obj["signature"].toString().toUtf8();
        return;
    }

    save(ml, res_json_obj["signature"].toString().toUtf8());
}

void LicenseManager::save(const QByteArray& l, const QByteArray& s) {
    qDebug() << "void LicenseManager::save";

    if (l.isEmpty() || s.isEmpty())
        return;

    qDebug() << "l.size()" << l.size();
    qDebug() << "s.size()" << s.size();

    QFile p(":/keys/private.pem");
    if (!p.open(QIODevice::ReadOnly)) {
        qDebug() << "!p.open(QIODevice::ReadOnly))";
        return;
    }
    QFile f(QString::fromUtf8(p.readAll()));
    p.close();

    if (f.open(QIODevice::WriteOnly)) {
        QDataStream out(&f);
        out << quint32(s.size());
        out.writeRawData(s.constData(), s.size());
        out << quint32(l.size());
        out.writeRawData(l.constData(), l.size());
        f.close();
        qDebug() << "OK";
    } else {
        qDebug() << "ERROR";
    }
}

void LicenseManager::load() {
    qDebug() << "void LicenseManager::load";

    QFile p(":/keys/private.pem");
    if (!p.open(QIODevice::ReadOnly)) {
        qDebug() << "!p.open(QIODevice::ReadOnly))";
        return;
    }
    QFile f(QString::fromUtf8(p.readAll()));
    p.close();
    if (!f.exists()) {
        qDebug() << "!f.exists()";
        remove();
        return;
    }

    if (f.open(QIODevice::ReadOnly)) {
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
        feed(b64pk[1].unicode(), 1);
        f.close();

        qDebug() << "l.size()" << l.size();
        qDebug() << "s.size()" << s.size();

        ml = l;
        ms = s;
    } else {
        remove();
        qDebug() << "!f.open(QIODevice::ReadOnly)";
    }
}

bool LicenseManager::issigned(const QByteArray& l, const QByteArray& s) {
    qDebug() << "issigned()";

    qDebug() << "l.toStdString()" << l.toStdString();
    qDebug() << "s.toStdString()" << s.toStdString();

    QFile f(":/keys/public.pem");
    if (!f.open(QIODevice::ReadOnly)) {
        qDebug() << "!f.open(QIODevice::ReadOnly))";
        return false;
    }
    QByteArray keyData = f.readAll();
    f.close();

    BIO* bio = BIO_new_mem_buf(keyData.data(), keyData.size());
    if (!bio)
        return false;

    EVP_PKEY* pkey = PEM_read_bio_PUBKEY(bio, nullptr, nullptr, nullptr);
    BIO_free(bio);

    if (!pkey) {
        qDebug() << "!pkey";
        return false;
    }

    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256(reinterpret_cast<const unsigned char*>(l.data()), l.size(), hash);

    QByteArray sig = QByteArray::fromBase64(s);

    EVP_PKEY_CTX* ctx = EVP_PKEY_CTX_new(pkey, nullptr);
    if (!ctx)
        return false;

    if (EVP_PKEY_verify_init(ctx) <= 0)
        return false;

    if (EVP_PKEY_CTX_set_rsa_padding(ctx, RSA_PKCS1_PADDING) <= 0)
        return false;

    if (EVP_PKEY_CTX_set_signature_md(ctx, EVP_sha256()) <= 0)
        return false;

    int res = EVP_PKEY_verify(ctx, reinterpret_cast<unsigned char*>(sig.data()), sig.size(), hash,
                              SHA256_DIGEST_LENGTH);

    EVP_PKEY_CTX_free(ctx);
    EVP_PKEY_free(pkey);

    qDebug() << "res == 1" << int(res == 1);

    return res == 1;
}

void LicenseManager::remove() {
    qDebug() << "void LicenseManager::remove";

    QFile p(":/keys/private.pem");
    if (!p.open(QIODevice::ReadOnly)) {
        qDebug() << "!p.open(QIODevice::ReadOnly))";
        return;
    }
    QFile f(QString::fromUtf8(p.readAll()));
    p.close();

    if (f.exists()) {
        if (!f.remove()) {
            qDebug() << "!f.remove()";
            return;
        }
    }
    ml.clear();
    ms.clear();

    h = 0xA5A5A5A5u;
    qDebug() << "OK";
}

void LicenseManager::license(const QByteArray& l) {
    qDebug() << "void LicenseManager::license";
    qDebug() << "l.toStdString()" << l.toStdString();
    ml = l;
    return;
}

// тут должен был быть самописный алгоритм, но не судьба
QByteArray LicenseManager::QEncodeByteArray(const QByteArray& l) const { return l.toBase64(); }

QByteArray LicenseManager::QDecodeByteArray(const QByteArray& l) const {
    return QByteArray::fromBase64(l);
}