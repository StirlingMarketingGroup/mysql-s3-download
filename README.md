# MySQL S3 Download

A small MySQL UDF library for downloading files from an AWS S3 bucket written in Golang.

---

## Usage

Know what you're doing! I'm sure this is wrong for a multitude reasons, but I'm personally using it for debug reasons, and because *sometimes* it's quicker to throw a query into Workbench then it is to write a script to view a file on S3, so, here it is.

---

### `s3_download`

Downloads a file from AWS S3 and returns the contents. Returns `NULL` if any of the params are `NULL`, or if the file doesn't exist, or there is some other error while downloading.

```sql
`s3_download` ( `region` , `bucket` , `key` )
```

- `` `region` ``
  - The AWS region that the bucket is in.
    - Example: `us-east-1`
- `` `bucket` ``
  - The AWS S3 bucket that the file is in.
- `` `key` ``
  - The key of the file you wish to download *without* the leading slash. Example: ""
    - Example: `path/to/file.csv`

## Examples

```sql
select`s3_download`('us-east-1','bucket.of.brian','path/to/file.csv');

--'Id,Name,Age
-- 1,Ahmad,21
-- 2,Ali,50'

select`s3_download`('us-east-1','bucket.of.brian','path/that/is/bad'); -- NULL
select`s3_download`('us-east-1',null,'path/to/file.csv'); -- NULL
```
---
## Downloading Errors

As you can see in the examples and the description above, the functions returns `NULL` if there is an error while downloading. If you do want to see the actual error messages, they will appear in the MySQL error log file like so:

```shell
tail -f /var/log/mysql/error.log # or wherever your MySQL error log is
s3-download: 2020/09/24 18:11:10.040207 /home/brian/go/src/github.com/StirlingMarketingGroup/mysql-s3-download/main.go:79: failed to download file from S3: NoSuchKey: The specified key does not exist.
        status code: 404, request id: ..., host id: ...
```

---

## Dependencies

You will need Golang, which you can get from here https://golang.org/doc/install. You will also need the MySQL dev library.

Debian / Ubuntu
```shell
sudo apt update
sudo apt install libmysqlclient-dev
```

## Configuring

Similar to other tools that connect to AWS, your AWS credentials must be set as environment variables. This is trivial if you're using something like Docker for running your database(s), but you can also set them on bare metal Ubuntu, for example, boxes by adding them to your `/etc/environment` file, or however your OS's environment variables are defined.

```shell
sudo nano /etc/environment

# Add these to the bottom

AWS_ACCESS_KEY_ID=Your key...
AWS_SECRET_ACCESS_KEY=Your secret...

```
## Installing

You can find your MySQL plugin directory by running this MySQL query

```sql
select @@plugin_dir;
```

then replace `/usr/lib/mysql/plugin` below with your MySQL plugin directory.

```shell
cd ~ # or wherever you store your git projects
git clone https://github.com/StirlingMarketingGroup/mysql-s3-download.git
cd mysql-s3-download
go build -buildmode=c-shared -o s3_download.so
sudo cp s3_download.so /usr/lib/mysql/plugin/ # replace plugin dir here if needed
```

Enable the functions in MySQL by running this MySQL query

```sql
create function`s3_download`returns string soname's3_download.so';
```