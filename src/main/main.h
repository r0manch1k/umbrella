#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QSoundEffect>

#include "./../redqueen/redqueen.h"

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE


class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:
    void spread();
    void about();
    void quit();

private:
    Ui::MainWindow *ui;
    RedQueenWidget *rq;
    QSoundEffect *click;
};

#endif // MAINWINDOW_H

