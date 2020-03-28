'''Executable module implementing Apache Beam pipeline.

This executable module implements Apache Beam pipeline importing patent
publications from Google Patents Public datasets.
'''

from __future__ import absolute_import

import argparse
import logging
import datetime

from apache_beam import Pipeline
from apache_beam import Map
from apache_beam import ParDo
from apache_beam.options.pipeline_options import PipelineOptions
from apache_beam.options.pipeline_options import SetupOptions
from apache_beam.io import ReadFromText
from apache_beam.transforms.combiners import ToList

from dofns import QueryPublications
from dofns import UpsertPatentsToDB

def run():
    '''Entry point, it defines and runs the pipeline.'''

    parser = argparse.ArgumentParser()

    parser.add_argument(
        '--db-user',
        dest='db_user',
        help='Username for Postgresql instance.',
        required=True
    )
    parser.add_argument(
        '--db-password',
        dest='db_password',
        help='Password for Postgresql instance.',
        required=True
    )
    parser.add_argument(
        '--db-name',
        dest='db_name',
        help='Patents database name.',
        required=True
    )
    parser.add_argument(
        '--db-host',
        dest='db_host',
        help='Hostname for Postgresql instance.',
        required=True
    )
    parser.add_argument(
        '--db-port',
        dest='db_port',
        help='Port number for Postgresql instance.',
        required=True
    )
    parser.add_argument(
        '--application-numbers-filepath',
        dest='application_numbers_filepath',
        help='Local or ``gs://`` path to file with application numbers.',
        required=True
    )

    known_args, pipeline_args = parser.parse_known_args()
    logging.info('Starting onboarding pipeline (args=%s)', known_args)

    # We use the `save_main_session` option because one or more DoFn's in this
    # workflow rely on global context (e.g., a module imported at module level).
    pipeline_options = PipelineOptions(pipeline_args)
    pipeline_options.view_as(SetupOptions).save_main_session = True

    with Pipeline(options=pipeline_options) as pipeline:
        # Read the text file[pattern] into a PCollection.
        lines = pipeline | 'ReadTextFile' >> ReadFromText(
            known_args.application_numbers_filepath
        )

        # Transform PCollection of text lines into PCollection of one element
        # which is slice of application numbers.
        app_numbers = lines | 'CombineApplicationNumbers' >> ToList()

        publications = (
            app_numbers
            | 'QueryPublications' >> ParDo(QueryPublications())
        )

        patents = (
            publications |
            'MapPublicationsToPatents' >> Map(
                lambda pb: {
                    'application_number': pb.application_number,
                    'application_kind': pb.application_kind,
                    'grant_date': datetime.datetime.fromtimestamp(pb.grant_date),
                },
            )
        )

        # Transform PCollection of patents into PCollection of one element which
        # is slice of patents.
        batch = patents | 'CombinePatentsToBatch' >> ToList()
        result = batch | 'UpsertPatentsToDB' >> ParDo(
            UpsertPatentsToDB(
                known_args.db_user,
                known_args.db_password,
                known_args.db_name,
                known_args.db_host,
                known_args.db_port,
            ),
        )

        return result

if __name__ == '__main__':
    logging.getLogger().setLevel(logging.INFO)
    run()
