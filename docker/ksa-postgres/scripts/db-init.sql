drop table if exists scanner;
create table scanner
(
    id serial
        constraint scanner_pk
            primary key,
    name varchar(32) not null
);
insert into scanner (name) values ('trivy'), ('prowler'), ('kube-bench');

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

alter table vulnerability add column search_vector tsvector;
create index vulnerability_search_vector_idx on vulnerability using gin (search_vector);

create function vulnerability_update_search_vector() returns trigger as $$
begin
    new.search_vector := to_tsvector('english', new.vid || ' ' || new.pkg_name || ' ' || new.title || ' ' || new.description || ' ' || new.target);
    return new;
end;
$$ language plpgsql;

create trigger vulnerability_search_vector_trigger
before insert or update on vulnerability
for each row execute function vulnerability_update_search_vector();

alter table misconfiguration add column search_vector tsvector;
create index misconfiguration_search_vector_idx on misconfiguration using gin (search_vector);

create function misconfiguration_update_search_vector() returns trigger as $$
begin
    new.search_vector := to_tsvector('english', new.mid || ' ' || new.type || ' ' || new.title || ' ' || new.description || ' ' || new.target);
    return new;
end;
$$ language plpgsql;

create trigger misconfiguration_search_vector_trigger
before insert or update on misconfiguration
for each row execute function misconfiguration_update_search_vector();
