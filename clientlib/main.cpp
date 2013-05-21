#include <QString>
#include "test.h"


int main()
{
    qDebug() << "Hi\n";
    TestSheerCloud test1;

    test1.SheerLinkUploadDownload();
    return 0;
}
