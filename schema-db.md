// ==========================================
// DATABASE 1: db_smartfactory_production
// ==========================================

Table production.user {
id uuid [primary key]
name varchar
email varchar [unique]
password varchar
role varchar
created_at timestamp
update_at timestamp
}

Table production.raw_materials {
id uuid [primary key]
material_name varchar
stock_quantity int
unit varchar
updated_at timestamp
}

Table production.machines {
id varchar [primary key]
name varchar
status varchar
created_at timestamp
updated_at timestamp
}

Table production.machine_alarm {
id uuid [primary key]
machine_id varchar [not null, ref: > production.machines.id] // Kita kunci relasinya langsung di sini (inline)
issue_description varchar  
 triggered_at timestamp
resolved_at timestamp [null]
}

// ==========================================
// DATABASE 2: db_smartfactory_iot
// ==========================================

Table iot.device_latest_metrics {
machine_id varchar [primary key]
temperature numeric
vibration numeric
last_seen timestamp
}

// ==========================================
// DATABASE 3: db_smartfactory_qc
// ==========================================

Table qc.quality_inspections {
id uuid [primary key]
machine_id varchar
image_path_url varchar
prediction_status varchar
confidence_score numeric
inspected_at timestamp
}

// ==========================================
// DATABASE 4: db_smartfactory_analytics (INFLUXDB - TSDB)
// ==========================================
// Catatan: Ini adalah database Time-Series NoSQL.
// Kotak ini merepresentasikan skema data streaming log jangka panjang.

Table influxdb.iot_telemetry_streams {
time timestamp [note: 'Automatic Index Time']
machine_id varchar [note: 'TAG - For High Speed Query Indexing']
temperature numeric [note: 'FIELD - Actual Value']
vibration numeric [note: 'FIELD - Actual Value']
}
