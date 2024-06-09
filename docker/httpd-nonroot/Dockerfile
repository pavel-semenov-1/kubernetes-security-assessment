FROM httpd:latest

RUN groupadd apache && \
    useradd -u 1000 -g 1000 apache && \
    chown -R apache:apache /usr/local/apache2 && \
    sed -i "s|Listen 80|Listen 8080|g" /usr/local/apache2/conf/httpd.conf 
    
USER 1000
EXPOSE 8080