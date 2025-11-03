#include "auth/auth.h"

#include <QApplication>
#include <QPalette>
#include <QMediaPlayer>
#include <QAudioOutput>
#include <QUrl>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    AuthWindow w;
    w.setFixedSize(1000, 600);
    w.show();
    return a.exec();
}
