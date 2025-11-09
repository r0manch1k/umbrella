#ifndef REDQUEENWIDGET_H
#define REDQUEENWIDGET_H

#include <QWidget>
#include <QStackedWidget>
#include <QVideoWidget>
#include <QMediaPlayer>
#include <QAudioOutput>
#include <QLabel>
#include <QGraphicsOpacityEffect>
#include <QPropertyAnimation>
#include <QAudioOutput>

class RedQueenWidget : public QWidget 
{
    Q_OBJECT

public:
    RedQueenWidget(QWidget *parent = nullptr);
    void enter();
    void about();
    void quit();

private:
    void say(const QString &path);

    QStackedWidget *stacked;
    QMediaPlayer *player;
    QVideoWidget *video;
    QLabel *image;
    QAudioOutput *audioOutput;
    QGraphicsOpacityEffect *photoEffect;

protected:
    void resizeEvent(QResizeEvent *event) override;
};
#endif // REDQUEENWIDGET_H