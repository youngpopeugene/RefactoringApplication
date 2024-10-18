create table if not exists substations
(
    name                 varchar(20) primary key,
    location             varchar(20) not null,
    year_of_construction integer     not null,
    commissioning_year   integer     not null
);
create table if not exists factories
(
    name varchar(20) primary key,
    city varchar(20) not null
);
create table if not exists range_of_high_voltage_equipments
(
    id                    serial primary key,
    high_voltage_switch   varchar(20) not null,
    medium_voltage_switch varchar(20) not null,
    low_voltage_switch    varchar(20) not null
);
create table if not exists cable_lines
(
    mark varchar(20) primary key
);
create table if not exists tire_sections
(
    name varchar(20) primary key
);

create table if not exists cell_kvls
(
    dispatch_name                  varchar(20) primary key,
    cable_line                     varchar(20) references Cable_lines,
    current_transformer            varchar(20) not null,
    switch                         varchar(20) not null,
    protection_transformer         varchar(20) not null,
    tire_section                   varchar(20) references Tire_sections,
    number_of_current_transformers integer     not null

);
create table if not exists fuses
(
    mark varchar(20) primary key
);
create table if not exists cell_tns
(
    dispatch_name       varchar(20) primary key,
    voltage_transformer varchar(20) not null,
    fuse                varchar(20) references Fuses,
    tire_section        varchar(20) references Tire_sections
);
create table if not exists cell_tsns
(
    dispatch_name         varchar(20) primary key,
    auxiliary_transformer varchar(20) not null,
    fuse                  varchar(20) references Fuses,
    tire_section          varchar(20) references Tire_sections
);
create table if not exists nsses
(
    id               serial primary key,
    rated_voltage_kV integer not null
);
create table if not exists range_of_standard_voltages
(
    id                          serial primary key,
    rated_winding_voltage_HV_kV integer not null,
    rated_winding_voltage_MV_kV integer not null,
    rated_winding_voltage_LV_kV integer not null
);
create table if not exists type_of_transformers
(
    type                      varchar(20) not null primary key,
    power_MVA                 integer     not null,
    cooling_system_type       varchar(20),
    range_of_standard_voltage serial references Range_of_standard_voltages
);
create table if not exists transformers
(
    factory_number                  integer primary key,
    NSS                             SERIAL references NSSes,
    substation                      varchar(20) references Substations,
    factory                         varchar(20) references Factories,
    type                            varchar(20) references Type_of_transformers,
    date_of_manufacture             date not null,
    commissioning_date              date not null,
    dispatch_name                   varchar(20),
    range_of_high_voltage_equipment serial references Range_of_high_voltage_equipments,
    tire_section                    varchar(20) references Tire_sections
);
create type roles AS enum (
    'WORKER',
    'DISPATCHER'
);
create table if not exists users
(
    username VARCHAR(20) primary key,
    password VARCHAR(100) not null,
    role roles not null
);
create table if not exists requests
(
    id                         serial primary key,
    transformer_factory_number integer references Transformers,
    worker_username            VARCHAR(20) references Users,
    is_completed               boolean default false,
    date_opened                date,
    date_closed                date

);