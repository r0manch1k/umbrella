#include "main.h"
#include "./ui_main.h"

#include <QPixmap>
#include <QSoundEffect>
#include <QPalette>
#include <QMediaPlayer>
#include <QAudioOutput>
#include <QUrl>
#include <QFontDatabase>
#include <QPushButton>
#include <QSoundEffect>

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent), ui(new Ui::MainWindow)
{
    ui->setupUi(this);

    QPixmap bg(":/images/background_main.jpg");
    bg = bg.scaled(this->size(), Qt::KeepAspectRatioByExpanding);

    QPalette palette;
    palette.setBrush(QPalette::Window, bg);
    this->setPalette(palette);

    QString family = QFontDatabase::applicationFontFamilies(QFontDatabase::addApplicationFont(":/fonts/SerpentineLight.ttf")).at(0);

    QFont fontU(family, 15);
    ui->titleTextUpperLabel->setFont(fontU);

    QFont fontL(family, 15);
    ui->titleTextLowerLabel->setFont(fontL);

    QFont fontA(family, 25);
    ui->rLabel_1->setFont(fontA);
    ui->aboutLabel_1->setFont(fontA);
    ui->aboutLabel_2->setFont(fontA);

    QVBoxLayout *layout = new QVBoxLayout(ui->rWidget);

    rq = new RedQueenWidget(this);
    layout->addWidget(rq);
    layout->setContentsMargins(0, 0, 0, 0);

    connect(rq, &RedQueenWidget::s_quit, this, [this]()
            { emit this->s_quit(); });

    connect(rq, &RedQueenWidget::s_spread, this, []()
            { QApplication::quit(); });

    click = new QSoundEffect(this);
    click->setSource(QUrl("qrc:/audio/click.wav"));
    click->setVolume(0.4);

    connect(ui->spreadButton, &QPushButton::clicked, this, &MainWindow::spread);
    connect(ui->aboutButton, &QPushButton::clicked, this, &MainWindow::about);
    connect(ui->quitButton, &QPushButton::clicked, this, &MainWindow::quit);
}

void MainWindow::spread()
{
    click->play();
    rq->spread();
}

void MainWindow::about()
{
    click->play();
    rq->about();
}

void MainWindow::quit()
{
    click->play();
    rq->quit();
}

MainWindow::~MainWindow()
{
    delete ui;
}
