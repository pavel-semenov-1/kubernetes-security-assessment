drop table if exists scanner;
create table scanner
(
    id serial
        constraint scanner_pk
            primary key,
    name varchar(32) not null
);
insert into scanner (name) values ('trivy');

drop table if exists report;
create table report
(
    id serial
        constraint report_pk
            primary key,
    scanner_id integer not null
        constraint report_scanner_fk
            references scanner (id),
    filename varchar(64) not null,
    parsed bool not null,
    generated_at timestamp not null
);

drop table if exists vulnerability;
create table vulnerability
(
    id serial
        constraint vulnerability_pk
            primary key,
    report_id integer not null
        constraint vulnerability_report_fk
            references report (id),
    vid varchar not null,
    pkg_name varchar not null,
    installed_version varchar not null,
    fixed_version varchar not null,
    title varchar not null,
    description varchar not null,
    severity varchar not null,
    target varchar not null
);

drop table if exists misconfiguration;
create table misconfiguration
(
    id serial
        constraint misconfiguration_pk
            primary key,
    report_id integer not null
        constraint misconfiguration_report_fk
            references report (id),
    mid varchar not null,
    type varchar not null,
    title varchar not null,
    description varchar not null,
    resolution varchar not null,
    severity varchar not null,
    target varchar not null
);