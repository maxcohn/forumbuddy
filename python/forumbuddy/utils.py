
class Validate:
    def numeric(value: str) -> int:
        '''Validate that an input is numeric and return it as an integer'''
        if value.isnumeric():
            return int(value)
        raise Exception('Value was not numeric')
        