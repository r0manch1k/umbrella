#include "redqueen.h"

#include <QStackedWidget>
#include <QMediaPlayer>
#include <QVideoWidget>
#include <QLabel>
#include <QAudioOutput>
#include <QVBoxLayout>
#include <QTimer>
#include <QUrl>
#include <QGraphicsOpacityEffect>
#include <QPropertyAnimation>

RedQueenWidget::RedQueenWidget(QWidget *parent)
    : QWidget(parent)
{
    image = new QLabel(this);
    QPixmap pixmap(":/images/redqueen.jpg");
    if (!pixmap.isNull()) {
        image->setPixmap(pixmap);
    }
    image->setAlignment(Qt::AlignCenter);
    image->setScaledContents(true);
    image->show();

    stacked = new QStackedWidget(this);
    stacked->setContentsMargins(0, 0, 0, 0);

    video = new QVideoWidget(stacked);
    video->setAspectRatioMode(Qt::IgnoreAspectRatio);
    stacked->addWidget(video);

    player = new QMediaPlayer(this);
    player->setVideoOutput(video);

    audioOutput = new QAudioOutput(this);
    audioOutput->setVolume(0.7);
    player->setAudioOutput(audioOutput);

    connect(player, &QMediaPlayer::mediaStatusChanged, this, [this](QMediaPlayer::MediaStatus status) {
        if (status == QMediaPlayer::EndOfMedia) {
            video->hide();
            player->stop();
        }
    });

    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->addWidget(stacked);
    layout->setContentsMargins(0,0,0,0);
    setLayout(layout);

    video->hide();

    QTimer::singleShot(0, this, &RedQueenWidget::enter);
}

void RedQueenWidget::enter() {
    say("resources/video/enter.mp4");
}

void RedQueenWidget::spread() {
    say("resources/video/spread.mp4");
}

void RedQueenWidget::about() {
    say("resources/video/about.mp4");
}

void RedQueenWidget::quit() {
    say("resources/video/quit.mp4");
}

void RedQueenWidget::say(const QString &path) {
    if (player->playbackState() == QMediaPlayer::PlayingState) {
        return;
    }
    video->raise();
    video->show();
    player->setSource(QUrl::fromLocalFile(path));
    player->setPosition(0);
    player->play();
}

void RedQueenWidget::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);
    stacked->setGeometry(0, 0, width(), height());
    image->setGeometry(0, 0, width(), height());
}