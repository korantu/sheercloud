#include "test.h"

#include <QByteArray>

void TestSheerCloudReally::VerifyTesting() {
  QVERIFY2(2+2==4, "Really");
}

void TestSheerCloudReally::SheerLinkLogin() {
  link.Authorize();
  loop.exec();
  QVERIFY2( link.Authorized(), "Password should match; Make sure the server is running.");
};

void TestSheerCloudReally::SheerLinkUploadDownload() {
  SheerLinkLogin();

  link.Upload("very/important/oldfile.txt", "123");
  loop.exec();

  QByteArray result;
  link.Download("very/important/oldfile.txt", result);
  loop.exec();

  QVERIFY2( result.contains("123"), "Sent/recieved data mismatch");
};

void TestSheerCloudReally::SheerLinkUploadDownloadBulk() {
  SheerLinkLogin();

  QByteArray massive("1234567890abcdefghijklmn"); // Every letter is a megabyte.
  massive = massive.repeated(1000000);
  link.Upload("very/huge/oldfile.txt", massive);
  loop.exec();

  QByteArray result;
  link.Download("very/huge/oldfile.txt", result);
  loop.exec();

  QVERIFY2( result == massive, "Sent/recieved data mismatch");
};

void TestSheerCloudReally::SheerLinkProgress() {
  SheerLinkLogin();

  QByteArray massive( "1234567890abcdefghijklmn"); // Every letter is a megabyte.
  massive = massive.repeated(1000000);
  int size = massive.size();

  link.Upload("very/huge/old_progressfile.txt", massive);
  int was = m_progress_reports;
  loop.exec();
  QVERIFY2( m_progress_reports - was > 3, "At least 3 reports expected");
  QVERIFY2( m_total == size, "Total is not correct");
  QVERIFY2( m_now == size, "Current is not correct");
  

  QByteArray result;
  link.Download("very/huge/old_progressfile.txt", result);
  loop.exec();

  QVERIFY2( result == massive, "Sent/recieved data mismatch");
};

void TestSheerCloudReally::SheerLinkDelete() {
  SheerLinkLogin();

  link.Upload("very/not_needed/file.txt", "123");
  loop.exec();

  link.Upload("very/not_needed/file_too.txt", "123");
  loop.exec();

  QByteArray result;
  link.Download("very/not_needed/file.txt", result);
  loop.exec();

  QVERIFY2( result.contains("123"), "Sent/recieved data mismatch");

  link.Delete("very/not_needed/file.txt");
  loop.exec();

  result.clear();
  link.Download("very/not_needed/file.txt", result);
  loop.exec();

  QVERIFY2( ! result.contains("123"), "Deleted file should have failed to be downloaded");
};

void TestSheerCloudReally::progress_check(qint64 now, qint64 total){
  m_now = now;
  m_total = total;
  m_progress_reports++;
}

QTEST_MAIN(TestSheerCloudReally)
