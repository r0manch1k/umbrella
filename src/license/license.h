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
    void issue(const QString &userId, int duration_hours, const QString &hwFingerprint);
    void verify();
    void save(const QByteArray &l, const QByteArray &s);
    void load();
    bool issigned(const QByteArray &l, const QByteArray &s);
    bool isLicenseValid() const;
    bool isExpired() const;
    QString license() const;

private:
    QByteArray m_l;
    QByteArray m_s;
    QDateTime m_exp;
    QNetworkAccessManager *nm;
    QByteArray loadPublicKey() const;
};
#endif // LICENSEMANAGER_H