#ifndef AUTHWINDOW_H
#define AUTHWINDOW_H

#include <QEvent>
#include <QMainWindow>
#include <QObject>
#include <QSoundEffect>

#include "../license/license.h"

QT_BEGIN_NAMESPACE
namespace Ui {
class AuthWindow;
}
QT_END_NAMESPACE

class AuthWindow : public QMainWindow {
    Q_OBJECT

   public:
    AuthWindow(QWidget* parent = nullptr);
    ~AuthWindow();

   private slots:
    void enter();

   private:
    void setLogoWidget();
    Ui::AuthWindow* ui;
    QSoundEffect* click;
    LicenseManager* lm;

   protected:
    bool eventFilter(QObject* obj, QEvent* event) override;
};

#endif  // AUTHWINDOW_H
