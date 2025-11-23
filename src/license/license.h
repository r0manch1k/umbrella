#ifndef LICENSEMANAGER_H
#define LICENSEMANAGER_H

#include <QAudioOutput>
#include <QGraphicsOpacityEffect>
#include <QJsonDocument>
#include <QJsonObject>
#include <QLabel>
#include <QMediaPlayer>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QPropertyAnimation>
#include <QStackedWidget>
#include <QVideoWidget>
#include <QWidget>

class LicenseManager : public QObject {
    Q_OBJECT

   public:
    explicit LicenseManager(QObject* parent = nullptr);
    void license(const QByteArray& l);
    void verify();
    void remove();
    QString fingerprint();
    void feed(uint32_t v, int n);
    void hash();
    void save(const QByteArray& l, const QByteArray& s);
    void load();
    uint32_t geth();
    QByteArray QEncodeByteArray(const QByteArray& l) const;
    QByteArray QDecodeByteArray(const QByteArray& l) const;
    bool issigned(const QByteArray& l, const QByteArray& s);
    bool isLicenseValid() const;
    bool isExpired() const;
    bool mv = false;
    QString b64pk;
    QByteArray ml;
    QByteArray ms;

   private:
    QDateTime me;
    QNetworkAccessManager* net;
    QByteArray loadPublicKey() const;
    uint32_t h = 0xA5A5A5A5u;
};
#endif  // LICENSEMANAGER_H