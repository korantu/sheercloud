#ifndef __test_h__
#define __test_h__

#include <sheer_cloud.h>

#include <QtTest/QtTest>

// Dunno where qint64 comes from.
#include <QString>
#include <QByteArray>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QFile>

class TestSheerCloudReally: public QObject {
  Q_OBJECT

  SheerCloudLink link;
  QEventLoop loop;

  qint64 m_now;
  qint64 m_total;
  int m_progress_reports;
 
 public:
 TestSheerCloudReally(): link("http://localhost:8080", "abc", "123") {
    link.connect( &link, SIGNAL(done()), &loop, SLOT(quit()));
    link.connect( &link, SIGNAL(progress(qint64 , qint64)), this, SLOT(progress_check(qint64 , qint64)));

   m_now = 0;
   m_total = 0;
   m_progress_reports = 0;

  };

 public slots:
  void progress_check(qint64 now, qint64 total);
    private slots:
  void VerifyTesting();
  void SheerLinkLogin();
  void SheerLinkUploadDownload();
  void SheerLinkUploadDownloadBulk();
  void SheerLinkProgress();
  void SheerLinkDelete();
};

#endif //  __test_h__
