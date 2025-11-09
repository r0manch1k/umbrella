#ifndef LICENSEMANAGER_H
#define LICENSEMANAGER_H

#include <QWidget>
#include <QStackedWidget>
#include <QVideoWidget>
#include <QMediaPlayer>
#include <QAudioOutput>
#include <QLabel>
#include <QGraphicsOpacityEffect>
#include <QPropertyAnimation>
#include <QAudioOutput>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QJsonObject>
#include <QJsonDocument>

class LicenseManager : public QObject
{
    Q_OBJECT

public:
    explicit LicenseManager(QObject *parent = nullptr);
    void requestLicenseFromServer(const QString &userId, int duration, const QString &hwFingerprint);
    void activateFromServerResponse(const QString &licenseJson, const QByteArray &signature);
    void verifyLicenseWithServer();
    bool verifySignature(const QByteArray &payload, const QByteArray &sig);
    void saveLicense(const QString &licenseJson);
    bool isLicenseValid() const;
    bool isExpired() const;
    QString currentLicense() const;

private:
    QByteArray m_licenseData;
    QByteArray m_signature;
    QDateTime m_expiration;
    bool m_valid = false;
    QNetworkAccessManager *networkManager;
    QByteArray loadPublicKey() const;
};
#endif // LICENSEMANAGER_H