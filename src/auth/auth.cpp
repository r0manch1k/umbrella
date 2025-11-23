#include "auth.h"

#include <QAudioOutput>
#include <QFontDatabase>
#include <QMediaPlayer>
#include <QMessageBox>
#include <QPalette>
#include <QPixmap>
#include <QSoundEffect>
#include <QTimer>
#include <QUrl>

#include "../license/license.h"
#include "./../main/main.h"
#include "./ui_auth.h"

AuthWindow::AuthWindow(QWidget* parent) : QMainWindow(parent), ui(new Ui::AuthWindow) {
    ui->setupUi(this);

    QPixmap bg(":/images/background.jpg");
    bg = bg.scaled(this->size(), Qt::KeepAspectRatioByExpanding);
    QPalette palette;
    palette.setBrush(QPalette::Window, bg);
    this->setPalette(palette);

    QString family = QFontDatabase::applicationFontFamilies(
                         QFontDatabase::addApplicationFont(":/fonts/SerpentineLight.ttf"))
                         .at(0);
    QFont fontU(family, 15);
    ui->titleTextUpperLabel->setFont(fontU);
    QFont fontL(family, 15);
    ui->titleTextLowerLabel->setFont(fontL);
    QFont fontN(family, 25);
    ui->nodeTextLabel->setFont(fontN);

    auto* audioOutput = new QAudioOutput(this);
    auto* player = new QMediaPlayer(this);
    player->setAudioOutput(audioOutput);
    player->setSource(QUrl("qrc:/audio/theme.mp3"));
    player->setLoops(QMediaPlayer::Infinite);
    audioOutput->setVolume(0.05);
    player->play();

    click = new QSoundEffect(this);
    click->setSource(QUrl("qrc:/audio/click.wav"));
    click->setVolume(0.4);

    ui->logoLabel->installEventFilter(this);

    connect(ui->enterButton, &QPushButton::clicked, this, &AuthWindow::enter);

    lm = new LicenseManager(this);

    lm->feed(lm->b64pk[0].unicode(), 0);
    lm->hash();

    lm->load();

    QTimer::singleShot(1000, this, [this]() {
        if (1628869386 * lm->issigned(lm->ml, lm->ms) == lm->geth()) {
            setLogoWidget();
        } else {
            lm->remove();
            lm->hash();
        }
    });
}

bool AuthWindow::eventFilter(QObject* obj, QEvent* event) {
    if (obj == ui->logoLabel && event->type() == QEvent::MouseButtonPress) {
        QMessageBox msgBox;
        msgBox.setText(lm->fingerprint());
        msgBox.exec();
        return true;
    }
    return QMainWindow::eventFilter(obj, event);
}

void AuthWindow::enter() {
    click->play();

    QString key = ui->keyLineEdit->text().trimmed();
    lm->hash();
    lm->license(key.toUtf8());
    lm->verify();
    lm->feed(lm->b64pk[0].unicode(), 0);
    lm->load();

    if (1628869386 * lm->issigned(lm->ml, lm->ms) == lm->geth()) {
        setLogoWidget();
    } else {
        lm->remove();
        ui->resLabel->setText("Access Denied");
    }
}

void AuthWindow::setLogoWidget() {
    lm->hash();
    if (lm->geth() != 826785547) {
        return;
    }
    auto* main = new MainWindow();
    connect(main, &MainWindow::s_quit, this, [this, main]() {
        this->move(main->pos());
        main->hide();
        this->show();
    });
    main->move(this->pos());
    main->show();
    this->hide();
}

AuthWindow::~AuthWindow() { delete ui; }
