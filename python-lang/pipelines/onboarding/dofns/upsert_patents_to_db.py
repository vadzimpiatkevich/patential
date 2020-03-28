'''Module providing definitions for upserting patents to database.

This module defines an Apache Beam DoFn class for upserting patents to database.
'''

from __future__ import absolute_import

import logging

from apache_beam import DoFn
from sqlalchemy.engine.url import URL
from sqlalchemy import create_engine
from sqlalchemy import MetaData
from sqlalchemy import Table
from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import Date
from sqlalchemy.dialects.postgresql import insert

POSTGRESQL_DRIVERNAME = 'postgresql+psycopg2'
TABLE_NAME = 'patents'

class UpsertPatentsToDB(DoFn):
    '''Apache Beam DoFn class for upserting patents to database.'''

    def __init__(self, db_user, db_password, db_name, db_host, db_port):
        super().__init__()

        self.db_uri = URL(
            drivername=POSTGRESQL_DRIVERNAME,
            username=db_user,
            password=db_password,
            database=db_name,
            host=db_host,
            port=db_port,
        )

        self.meta = None
        self.connection = None

    def setup(self):
        engine = create_engine(self.db_uri)
        self.meta = MetaData(engine)

        logging.info('Connecting to Postgresql (uri=%s)', self.db_uri)
        self.connection = engine.connect()

    def process(self, batch):
        table = Table(
            TABLE_NAME,
            self.meta,
            Column('application_number', String, unique=True, nullable=False),
            Column('application_kind', String, nullable=False),
            Column('grant_date', Date, nullable=False),
        )

        stmt = insert(table).values(batch)
        upsert_stmt = stmt.on_conflict_do_nothing(
            index_elements=['application_number']
        )

        logging.info('Upserting batch to Postgresql (size=%d)', len(batch))
        return self.connection.execute(upsert_stmt)

    def teardown(self):
        logging.info('Closing connection to Postgresql')
        self.connection.close()
