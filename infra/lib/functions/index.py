import json
import requests
from opentelemetry import trace, baggage

tracer = trace.get_tracer(__name__)


def handler(event, context):
    print(json.dumps(event))

    current_span = trace.get_current_span()
    current_span.set_attribute("http.route", "some_route")

    with tracer.start_as_current_span("server_span") as span:
        requests.get("https://aws.amazon.com/")

    with tracer.start_as_current_span("error span") as span:
        span.set_attribute("projectID", baggage.get_baggage("projectID"))
        span.add_event("event message", {"event_attributes": 1})

        try:
            1 / 0
        except ZeroDivisionError as err:
            span.record_exception(err)
            print("caught zero division error")

    return "ok"