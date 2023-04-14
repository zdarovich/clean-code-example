create database if not exists test_database;
use test_database;
create table if not exists todo (id varchar(50),title varchar(255),completed boolean,created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
