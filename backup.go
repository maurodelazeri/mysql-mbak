package main

import (
    "github.com/keighl/barkup"
    "fmt"
    "strconv"
    "time"
)

func RunBackupProcess() error {
    for _, host := range config.Connections {
        logger.Info("starting backup for host %s:%d", host.Hostname, host.Port)
        for _, db := range host.Databases {
            if err := RunBackup(host, db); err != nil {
                return err
            }
        }
    }

    return nil
}

func RunBackup(host ConnectionConfig, dbName string) error {

    logger.Info("| %s", dbName)

    mysql := BuildMysqlConfig(host, dbName)

    if host.S3Bucket != "" && host.S3Path != "" {

        s3 := BuildS3Config(host.S3Bucket)

        t := time.Now().UTC()

        path := fmt.Sprintf("%s%d/%s/%d/", host.S3Path, t.Year(), t.Month().String(), t.Day())

        logger.Debug("backing up to s3 bucket: %s", path)
        if err := mysql.Export().To(path, s3); err != nil {
            return fmt.Errorf("failed to run backup for %s:%d/%s\n%s", host.Hostname, host.Port, dbName, err.Error())
        }

    } else if host.LocalDir != "" {

        t := time.Now().UTC()

        path := fmt.Sprintf("%s%d/%s/%d/", host.LocalDir, t.Year(), t.Month().String(), t.Day())

        logger.Debug("backing up to local directory: %s", path)
        if err := mysql.Export().To(host.LocalDir, nil); err != nil {
            return fmt.Errorf("failed to run backup for %s:%d/%s\n%s", host.Hostname, host.Port, dbName, err.Error())
        }

    } else {

        return fmt.Errorf("no storage destination specified for %s:%d/%s. Must specify (local_directory) or (s3_bucket and s3_path)", host.Hostname, host.Port, dbName)

    }

    return nil
}

func BuildMysqlConfig(host ConnectionConfig, dbName string) *barkup.MySQL {
    logger.Debug("building MySQL config")
    return &barkup.MySQL{
        Host:       host.Hostname,
        Port:       strconv.Itoa(host.Port),
        User:       host.Username,
        Password:   host.Password,
        DB:         dbName,
    }
}

func BuildS3Config(bucket string) *barkup.S3 {
    logger.Debug("building s3 config")

    return &barkup.S3{
        Region:         config.S3Config.Region,
        AccessKey:      config.S3Config.AccessKey,
        ClientSecret:   config.S3Config.ClientSecret,
        Bucket:         bucket,
    }
}
