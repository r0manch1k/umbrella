#include <QApplication>
#include <QAudioOutput>
#include <QIcon>
#include <QMediaPlayer>
#include <QPainter>
#include <QPalette>
#include <QUrl>

#include "auth/auth.h"

QIcon icofixed(const QPixmap& px, int pad) {
    QPixmap out(px.width() + pad * 2.5, px.height() + pad * 2.5);
    out.fill(Qt::transparent);
    QPainter p(&out);
    p.drawPixmap(pad, pad, px);
    return QIcon(out);
}

int main(int argc, char* argv[]) {
    QApplication a(argc, argv);
    a.setOrganizationName("Umbrella Corp.");
    a.setOrganizationDomain("theosokolovsky.ru");
    a.setApplicationName("Umbrella Corp. INAN");
    a.setApplicationVersion("1.0.0");
    a.setWindowIcon(icofixed(QPixmap(":/images/logo.png"), 12));
    AuthWindow w;
    w.setFixedSize(1000, 600);
    w.show();
    return a.exec();
}
