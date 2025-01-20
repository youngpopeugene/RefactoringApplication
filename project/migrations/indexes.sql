CREATE INDEX idx_substation_location ON substations (location);

CREATE INDEX idx_transformer_substation ON transformers (substation);

CREATE INDEX idx_transformer_factory_number ON transformers (factory_number);
