#include "auth.h"
#include "./ui_auth.h"

#include <QPixmap>
#include <QPalette>

#include <QMediaPlayer>
#include <QAudioOutput>
#include <QUrl>
#include <QFrame>
#include <QFontDatabase>
#include <QLabel>
#include <QShortcut>
#include <QTextEdit>
#include <QMouseEvent>
#include <QWidget>
#include <QMainWindow>
#include <QApplication>



AuthWindow::AuthWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::AuthWindow)
{
    ui->setupUi(this);

    int id = QFontDatabase::addApplicationFont(":/fonts/SerpentineLight.ttf");
    QString family = QFontDatabase::applicationFontFamilies(id).at(0);

    QFont fontU(family, 15);
    ui->titleTextUpperLabel->setFont(fontU);

    QFont fontL(family, 15);
    ui->titleTextLowerLabel->setFont(fontL);

    QFont fontN(family, 25);
    ui->nodeTextLabel->setFont(fontN);

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

    QFrame *overlay = new QFrame(this);
    overlay->setStyleSheet("background-color: rgba(0, 0, 0, 180);");
    overlay->setGeometry(0, 0, 300, 600);
    overlay->hide();

    QTextEdit *console = new QTextEdit(overlay);
    console->setGeometry(0, 0, 300, 600);
    console->setStyleSheet("background-color: black; color: lime;");


    QShortcut *shortcut = new QShortcut(QKeySequence("~"), this);
    connect(shortcut, &QShortcut::activated, this, [=]() {
        overlay->setVisible(!overlay->isVisible());
    });
}

void AuthWindow::mousePressEvent(QMouseEvent *event)
{
    QWidget *focused = QApplication::focusWidget();
    if (focused && qobject_cast<QLineEdit*>(focused)) {
        focused->clearFocus();
    }
    QMainWindow::mousePressEvent(event);
}

AuthWindow::~AuthWindow()
{
    delete ui;
}

