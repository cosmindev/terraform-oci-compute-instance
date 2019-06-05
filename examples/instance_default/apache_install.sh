#
# @Author: cosmin.tudor@oracle.com 
# @Date: 2018-11-08 12:37:09 
# @Last Modified by:   cosmin.tudor@oracle.com 
# @Last Modified time: 2018-11-08 12:37:09 
#

#!/bin/bash -x

INSTANCE_DNS=$1

# Install apache
yum install -y httpd

echo 'Hello World from Apache running on '${INSTANCE_DNS}'' > /var/www/html/index.html

# echo "Listen 80" >> /etc/httpd/conf/httpd.conf
service httpd start

# make httpd service start at boot
chkconfig --add httpd
chkconfig httpd on

#enable port 80 at the OS firewall level
sudo firewall-cmd --permanent --zone=public --add-service=http 
sudo firewall-cmd --reload

# END install apache
