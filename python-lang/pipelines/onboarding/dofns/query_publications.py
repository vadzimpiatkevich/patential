'''Module providing definitions for querying patents publications.

This module defines an Apache Beam DoFn class for querying patent publications
from BigQury by application numbers.
'''

from __future__ import absolute_import

import logging

from apache_beam import DoFn
from google.cloud.bigquery import Client
from google.cloud.bigquery import QueryJobConfig
from google.cloud.bigquery import ScalarQueryParameter
from google.cloud.bigquery import ArrayQueryParameter

US_COUNTRY_CODE = 'US'
WIPO_KIND_CODES_FROM = 20010102

USPTO_PATENT_PUBLICATION_CODE = 'A'
WIPO_PATENT_PUBLICATION_CODES = ['B1', 'B2']

class QueryPublications(DoFn):
    '''Apache Beam DoFn class for querying patent publications from BigQuery by
    application numbers.
    '''

    def __init__(self):
        super().__init__()
        self.storage_client = None

    def setup(self):
        logging.info('Initializing BigQuery client')
        self.storage_client = Client()

    def process(self, app_numbers):
        sql = """
            SELECT application_number, application_kind, grant_date
            FROM `patents-public-data.patents.publications`
            WHERE
            country_code = @us_country_code
            AND application_number IN UNNEST(@application_numbers)
            AND IF (
                publication_date >= @wipo_kind_codes_from,
                kind_code IN UNNEST(@wipo_patent_publication_codes),
                kind_code = @uspto_patent_publication_code
            );
        """

        job_config = QueryJobConfig(
            query_parameters=[
                ScalarQueryParameter(
                    'us_country_code',
                    'STRING',
                    US_COUNTRY_CODE,
                ),
                ArrayQueryParameter(
                    'application_numbers',
                    'STRING',
                    app_numbers,
                ),
                ScalarQueryParameter(
                    'wipo_kind_codes_from',
                    'INT64',
                    WIPO_KIND_CODES_FROM,
                ),
                ArrayQueryParameter(
                    'wipo_patent_publication_codes',
                    'STRING',
                    WIPO_PATENT_PUBLICATION_CODES,
                ),
                ScalarQueryParameter(
                    'uspto_patent_publication_code',
                    'STRING',
                    USPTO_PATENT_PUBLICATION_CODE,
                ),
            ]
        )
        query = self.storage_client.query(sql, job_config=job_config)

        logging.info('Executing query for publications')
        iterator = query.result()

        return iterator
