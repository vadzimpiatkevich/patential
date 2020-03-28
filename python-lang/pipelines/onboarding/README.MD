Apache Beam pipeline for importing patent publications from Google Patents Public datasets.

## Overview:

This pipeline connects to Bigquery in order to retrieve grant-patent publications by specified application numbers. It maps publications to Patential data model and upserts them to destination Patents Postgresql database.

## Local development:

Before you begin:

1. Setup python development environment: https://cloud.google.com/python/setup

2. Install pipeline dependencies: `pip install -r requirements.txt`

3. Setup Postgresql development instance.
Quickstart with local instance as an example: https://www.codementor.io/@engineerapart/getting-started-with-postgresql-on-mac-osx-are8jcopb

4. Guarantee python client libraries have access to BigQuery and PubSub API.
Quickstart: https://cloud.google.com/bigquery/docs/quickstarts/quickstart-client-libraries

The pipeline requires multiple command line arguments:

* --db-user: Username for Postgresql instance.
* --db-password: Password for Postgresql instance.
* --db-name: Patents database name.
* --db-host: Hostname for Postgresql instance.
* --db-port: Port number for Postgresql instance.
* --application-numbers-filepath: Local or ``gs://`` path to file with application numbers.

Example:
```
python main.py \
--db-user postgres \
--db-password password \
--db-name patential_patents \
--db-host localhost \
--db-port 5432 \
--application-numbers-filepath application-numbers-sample.txt
```
