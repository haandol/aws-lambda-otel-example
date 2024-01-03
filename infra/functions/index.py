import json
import requests
import traceback
from opentelemetry import trace
from opentelemetry.semconv.trace import SpanAttributes


tracer = trace.get_tracer(__name__)


def handler(event, context):
    print(json.dumps(event))

    with tracer.start_as_current_span("requests amazon") as span:
        url = "https://aws.amazon.com/"
        requests.get(
            url,
            {
                "timeout": 1,
            },
        )
        span.set_attribute(SpanAttributes.HTTP_METHOD, "GET")
        span.set_attribute(SpanAttributes.HTTP_URL, url)

    with tracer.start_as_current_span("devide zero") as span:
        try:
            1 / 0
        except ZeroDivisionError as err:
            span.set_status(trace.StatusCode.ERROR)
            span.record_exception(
                err,
                {
                    SpanAttributes.EXCEPTION_STACKTRACE: traceback.format_exc(),
                    SpanAttributes.EXCEPTION_ESCAPED: True,
                },
            )
            return {
                "statusCode": 500,
                "body": json.dumps(
                    {
                        "stacktrace": traceback.format_exc(),
                        "message": "zero division error",
                    }
                ),
            }

    return {
        "statusCode": 200,
        "body": json.dumps(
            {
                "message": "ok",
            }
        ),
    }
