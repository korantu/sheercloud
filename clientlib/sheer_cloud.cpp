#include "sheer_cloud.h"

#include <QDebug>

SheerCloudLink::SheerCloudLink(QString location, QString login, QString password){
  m_is_authorized = false;

  m_location = location;
  m_login = login;
  m_password = password;
  m_out = NULL;
};

SheerCloudLink::~SheerCloudLink(){
};
  
void SheerCloudLink::Authorize(){
  m_reply = get( QNetworkRequest( QUrl( m_location + "/authorize?login=" + m_login + "&password=" + m_password ) ));
  connect( m_reply, SIGNAL(finished()), this, SLOT(login_completed()) );
}

void SheerCloudLink::Upload(QString file, const QByteArray & in){
  QNetworkRequest upload_req( QUrl( m_location + "/upload?login=" + m_login + "&password=" + m_password + "&file=" + file ));
  upload_req.setRawHeader( "content-type", "application/octet-stream");
  m_reply = post( upload_req, in);
  connect(m_reply, SIGNAL(finished()), SLOT(upload_completed()));
  connect(m_reply, SIGNAL(uploadProgress ( qint64 , qint64 ) ), SIGNAL(progress ( qint64 , qint64 ) ));

};

void SheerCloudLink::Download(QString file, QByteArray & out){
  m_reply = get( QNetworkRequest( QUrl( m_location + "/download?login=" + m_login + "&password=" + m_password + "&file=" + file ) ));
  m_out = &out;
  connect(m_reply, SIGNAL(finished()), SLOT(download_completed()));
  connect(m_reply, SIGNAL(downloadProgress ( qint64 , qint64 ) ), this, SIGNAL(progress ( qint64 , qint64 ) ));
};

void SheerCloudLink::List(QString file, QByteArray & out){
  m_reply = get( QNetworkRequest( QUrl( m_location + "/list?login=" + m_login + "&password=" + m_password + "&file=" + file ) ));
  m_out = &out;
  connect(m_reply, SIGNAL(finished()), SLOT(download_completed()));
  connect(m_reply, SIGNAL(downloadProgress ( qint64 , qint64 ) ), SIGNAL(progress ( qint64 , qint64 ) ));
};

QList<CloudFile> ParseList( const QByteArray & in){
  QStringList list = QString(in).split("\n", QString::SkipEmptyParts);
  // TODO iterate over list, how it is done in QT
  QList<CloudFile> result;
  for ( int i = 0; i+1 < list.size(); i+=2){
    CloudFile entry;
    entry.name = list.at(i);
    entry.hash = list.at(i+1);
    result.push_back(entry);
  }
  return result;
}

void SheerCloudLink::Delete(QString file){
  QNetworkRequest upload_req( QUrl( m_location + "/delete?login=" + m_login + "&password=" + m_password + "&file=" + file ));
  m_reply = get( upload_req);
  connect( m_reply, SIGNAL(finished() ), this, SLOT(delete_completed()) );
};

void SheerCloudLink::Job(QString file, JobID &out){
  job_id = &out;
  QNetworkRequest upload_req( QUrl( m_location + "/job?login=" + m_login + "&password=" + m_password + "&file=" + file ));
  m_reply = get( upload_req);
  connect( m_reply, SIGNAL(finished() ), this, SLOT(job_requested()) );
};

void SheerCloudLink::Progress(JobID id, JobResult &out){
  job_result = &out;
  QNetworkRequest upload_req( QUrl( m_location + "/progress?login=" + m_login + "&password=" + m_password + "&id=" + id ));
  m_reply = get( upload_req);
  connect( m_reply, SIGNAL(finished() ), this, SLOT(progress_requested()) );
};

bool SheerCloudLink::Authorized(){
  return m_is_authorized;
};

void SheerCloudLink::request_completed(){
  disconnect( this, SLOT(login_completed()) );
  disconnect( this, SLOT(download_completed()) );
  disconnect( this, SLOT(upload_completed()) );
  disconnect( this, SLOT(delete_completed()) );
  disconnect( this, SIGNAL(progress(qint64, qint64)));
  m_reply->deleteLater();
  m_reply = NULL;
  done();
}

void SheerCloudLink::login_completed(){
  if( QString(m_reply->readAll()).contains( "OK" ) ) {
    m_is_authorized = true;
  };
  request_completed();
};

void SheerCloudLink::upload_completed(){
  QByteArray got = m_reply->readAll();
  request_completed();
};

void SheerCloudLink::download_completed(){
  QByteArray got = m_reply->readAll();
  if( m_out != NULL ) {
    *m_out = got;
  };
  request_completed();
};

void SheerCloudLink::delete_completed(){
  QByteArray got = m_reply->readAll();
  request_completed();
};


void SheerCloudLink::job_requested(){
  QByteArray out = m_reply->readAll();
  *job_id = QString(out).replace("OK:", "");
  request_completed();
};


void SheerCloudLink::progress_requested(){
  QByteArray got = m_reply->readAll();
  // qDebug() << got;
  *job_result = ( QString(got) == "OK:DONE" );
  request_completed();
};
