// generated by tabler
package models

// User

func (u User) CreateTable() string {
    return `CREATE TABLE user (id uuid, email varchar(128), created timestamp) PRIMARY KEY (id);`
}

func (u User) DropTable() string {
    return `DROP TABLE user;`
}

func (u User) InsertRow() string {
    return `INSERT INTO user (id, email, created) VALUES (?, ?, ?);`
}

func (u User) SelectRow() string {
    return `SELECT id, email, created FROM user WHERE id=?;`
}

// Profile

func (p Profile) CreateTable() string {
    return `CREATE TABLE profile (userid uuid REFERENCES user(id), attribute varchar(64), value varchar(256)) PRIMARY KEY (userid, attribute);`
}

func (p Profile) DropTable() string {
    return `DROP TABLE profile;`
}

func (p Profile) InsertRow() string {
    return `INSERT INTO profile (userid, attribute, value) VALUES (?, ?, ?);`
}

func (p Profile) SelectRow() string {
    return `SELECT userid, attribute, value FROM profile WHERE userid=? AND attribute=?;`
}

// NoPrimary

func (n NoPrimary) CreateTable() string {
    return `CREATE TABLE noprimary (attribute varchar(64), value varchar(256));`
}

func (n NoPrimary) DropTable() string {
    return `DROP TABLE noprimary;`
}

func (n NoPrimary) InsertRow() string {
    return `INSERT INTO noprimary (attribute, value) VALUES (?, ?);`
}

func (n NoPrimary) SelectRow() string {
    return `SELECT attribute, value FROM noprimary WHERE ;`
}
