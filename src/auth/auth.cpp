#include "auth.h"
#include "./ui_auth.h"
#include "./../main/main.h"
#include "../license/license.h"

#include <QPixmap>
#include <QSoundEffect>
#include <QPalette>
#include <QMediaPlayer>
#include <QAudioOutput>
#include <QUrl>
#include <QFontDatabase>
#include <QMessageBox>
#include <QTimer>


AuthWindow::AuthWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::AuthWindow)
{
    ui->setupUi(this);

    QPixmap bg(":/images/background.jpg");
    bg = bg.scaled(this->size(), Qt::KeepAspectRatioByExpanding);

    QPalette palette;
    palette.setBrush(QPalette::Window, bg);
    this->setPalette(palette);

    QString family = QFontDatabase::applicationFontFamilies(QFontDatabase::addApplicationFont(":/fonts/SerpentineLight.ttf")).at(0);

    QFont fontU(family, 15);
    ui->titleTextUpperLabel->setFont(fontU);

    QFont fontL(family, 15);
    ui->titleTextLowerLabel->setFont(fontL);

    QFont fontN(family, 25);
    ui->nodeTextLabel->setFont(fontN);

    auto *audioOutput = new QAudioOutput(this);
    auto *player = new QMediaPlayer(this);
    player->setAudioOutput(audioOutput);
    player->setSource(QUrl("qrc:/audio/theme.mp3"));
    player->setLoops(QMediaPlayer::Infinite);
    audioOutput->setVolume(0.1);
    player->play();

    click = new QSoundEffect(this);
    click->setSource(QUrl("qrc:/audio/click.wav"));
    click->setVolume(0.4);

    ui->logoLabel->installEventFilter(this);
    
    connect(ui->enterButton, &QPushButton::clicked, this, &AuthWindow::enter);

    licenseManager = new LicenseManager(this);
}

bool AuthWindow::eventFilter(QObject *obj, QEvent *event)
{
    if (obj == ui->logoLabel && event->type() == QEvent::MouseButtonPress) {
        buy();
        return true;
    }
    return QMainWindow::eventFilter(obj, event);
}


void AuthWindow::enter()
{
    click->play();

    QString key = ui->keyLineEdit->text().trimmed();

    QByteArray signature;

    licenseManager->verify();

    // if (licenseManager->isLicenseValid()) {
    if (1) {
        auto *main = new MainWindow();
        main->move(this->pos());
        main->show();
        this->hide();
    } else {
        ui->resLabel->setText("Access Denied");
    }
}

void AuthWindow::buy()
{
    click->play();

    QString userId = "user123";
    int duration = 30;
    QString hwFingerprint = "HW1234";

    licenseManager->issue(userId, duration, hwFingerprint);

    QMessageBox::information(this, "ACCESS KEY", licenseManager->license());
}


AuthWindow::~AuthWindow()
{
    delete ui;
}
