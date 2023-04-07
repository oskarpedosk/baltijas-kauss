#!/bin/bash

# Create new user and add to sudo group
sudo adduser op
sudo usermod -aG sudo op

# Update package lists and upgrade existing packages
sudo apt update
sudo apt upgrade

# Install PostgreSQL 14 and check its status
sudo apt install postgresql-14
service postgresql status

# Install required packages for Caddy web server
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy

# Install Supervisor process control system
sudo apt install supervisor

# Configure SSH
sudo nano /etc/ssh/sshd_config
sudo systemctl restart sshd

# Install Go 1.20.3
wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile

# Clone Git repository
git clone https://github.com/oskarpedosk/baltijas-kauss.git
cd baltijas-kauss
git config --global user.email "oskar.pedosk@gmail.com"
git config --global user.name "Oskar Pedosk"
git config --global credential.helper store

# THIS PART DOESNT WORK 
echo "#ACCESS TOKEN" | git credential-store --file ~/.git-credentials store
# THIS PART DOESNT WORK 

# Edit PostgreSQL configuration file
# sudo nano /etc/postgresql/14/main/pg_hba.conf
# Replace "peer" with "trust" in the "pg_hba.conf" file for "localhost" and "127.0.0.1" entries
sudo sed -i 's/host[[:space:]]*all[[:space:]]*all[[:space:]]*127\.0\.0\.1\/32[[:space:]]*scram-sha-256/host    all             all             127.0.0.1\/32            trust/g' /etc/postgresql/14/main/pg_hba.conf
sudo sed -i 's/host[[:space:]]*all[[:space:]]*all[[:space:]]*::1\/128[[:space:]]*scram-sha-256/host    all             all             ::1\/128                 trust/g' /etc/postgresql/14/main/pg_hba.conf

# Restart PostgreSQL service and check its status
sudo service postgresql stop
sudo service postgresql start
ps ax | grep postgr


# Connect DBeaver HERE
# Connect DBeaver use SSH tunnel
# Select show all databases (edit connection)

cd
cd baltijas-kauss
cp database.yml.example database.yml
# nano database.yml
# user: postgres
sed -i "s/user: /user: postgres/g" database.yml
go get github.com/gobuffalo/pop/...
go install github.com/gobuffalo/pop/v6/soda@latest
ls ~/go/bin
cd
# nano .profile
echo 'export PATH=$PATH:~/go/bin' >> ~/.profile
export PATH=$PATH:~/go/bin
cd baltijas-kauss
soda migrate
go build -o baltijaskauss cmd/web/*.go
ls -l
cd /etc/caddy
sudo mv Caddyfile Caddyfile.dist
sudo nano Caddyfile

# Add the specified configuration to the file
sudo sh -c "echo '{
	email oskar.pedosk@gmail.com
}
(static) {
	@static {
		file
		path *.ico *.css *.js *.mp3 *.gif *.jpg *.jpeg *.png *.svg *.webp *.woff *.json
	}
	header @static Cache-Control max-age=5184000
}	

(security) {
    header {
        # Enable HSTS
		Strict-Transport-Security max-age=31536000;
		# Disable clients from sniffing media type
		X-Content-Type-Options nosniff
		# Keep referrer data off of HTTP connections
		Referrer-Policy no-referrer-when-downgrade
    }
}

import conf.d/*.conf
' >> /etc/caddy/Caddyfile"

# Create and edit a new configuration file in conf.d/
sudo mkdir /etc/caddy/conf.d

sudo sh -c "echo '
baltijaskauss.ee {
    encode zstd gzip
    import static
    import security

    log {	
	    output file /var/www/bkauss/logs/caddy.access.log
    }	
	
    reverse_proxy http://localhost:8080
}
' >> /etc/caddy/conf.d/bkauss.conf"

cd /var
sudo mkdir www
cd www
sudo mv ~/baltijas-kauss bkauss
cd bkauss
# mkdir logs
sudo chmod 777 logs
sudo service caddy restart
sudo service caddy status
cd /var/www/bkauss
./baltijaskauss -h
./baltijaskauss -dbname=baltijas_kauss -dbpass=bkauss1 -dbuser=postgres

cd /etc/supervisor
cd conf.d
sudo nano bkauss.conf
[program:bkauss]
command=/var/www/bkauss/baltijaskauss -dbname=baltijas_kauss -dbpass=bkauss1 -dbuser=postgres
autorestart=true
autostart=true
stdout_logfile=/var/www/bkauss/logs/supervisord.log

sudo sh -c "echo '[program:bkauss]
command=/var/www/bkauss/baltijaskauss -dbname=baltijas_kauss -dbpass=bkauss1 -dbuser=postgres
directory=/var/www/bkauss
autorestart=true
autostart=true
stdout_logfile=/var/www/bkauss/logs/supervisord.log
' >> bkauss.conf"
sudo supervisorctl
status
update
status

# Manually
# Create update.sh
cd /var/www/bkauss
sudo sh -c "echo '#!/bin/bash
git reset --hard
git pull
soda migrate
go build -o baltijaskauss cmd/web/*.go
sudo supervisorctl stop bkauss
sudo supervisorctl start bkauss
' >> update.sh"

sudo chmod 777 update.sh

# Run update.sh
./update.sh << EOL
oskar.pedosk@gmail.com
ACCESS TOKEN
EOL

# Update package list
sudo apt-get update

# Install node.js and npm
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs
sudo apt-get install -y npm

# Install chromium browser
sudo apt-get install -y chromium-browser

# Install puppeteer and stealth plugin
sudo npm install -g puppeteer puppeteer-extra puppeteer-extra-plugin-stealth
