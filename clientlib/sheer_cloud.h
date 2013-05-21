#ifndef __cloud_h__
#define __cloud_h__

#include <QString>
#include <QByteArray>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QFile>


class SheerCloudLink: public QNetworkAccessManager {
  Q_OBJECT;
 private:
  QString m_location, m_login, m_password;

  bool m_is_authorized;
  QNetworkReply * m_reply; // Current request

  QByteArray * m_out; // Keep track of output

  void request_completed(); // Prepare for the next request

 public:
  SheerCloudLink(QString location, QString login, QString password);
  virtual ~SheerCloudLink();
  
  void Authorize();
  bool Authorized();

  void Upload(QString, const QByteArray &);
  void Download(QString, QByteArray &);
  void Delete(QString);

  signals:
  void done();
  void progress(qint64 bytesReceived, qint64 bytesTotal);

 private slots:
  void login_completed();
  void upload_completed();
  void download_completed();
  void delete_completed();
};

#endif
