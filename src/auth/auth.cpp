#include "auth.h"
#include "./ui_auth.h"

#include <QPixmap>
#include <QPalette>

#include <QMediaPlayer>
#include <QAudioOutput>
#include <QUrl>


AuthWindow::AuthWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::AuthWindow)
{
    ui->setupUi(this);

    QPixmap bkgnd(":/images/background.jpg");
    bkgnd = bkgnd.scaled(this->size(), Qt::KeepAspectRatioByExpanding);

    QPalette palette;
    palette.setBrush(QPalette::Window, bkgnd);
    this->setPalette(palette);

    auto *audioOutput = new QAudioOutput(this);
    auto *player = new QMediaPlayer(this);
    player->setAudioOutput(audioOutput);

    player->setSource(QUrl("qrc:/audio/theme.mp3"));
    player->setLoops(QMediaPlayer::Infinite);
    audioOutput->setVolume(0.25);

    player->play();
}

AuthWindow::~AuthWindow()
{
    delete ui;
}

