"""Helping the bridge engineers with Python."""

import dataclasses

_DATA_PATH = '../data/day7.txt'


@dataclasses.dataclass
class Equation:
    value: int
    inputs: list[int]


def parseLine(line: str) -> Equation:
    """Parse a calibration equation."""
    value, inputs = line.split(':', 1)
    value = int(value)
    inputs = [int(i) for i in inputs.strip().split(' ')]
    return Equation(value=value, inputs=inputs)


def parse(data: str) -> list[Equation]:
    """Parse a set of calibration equations."""
    return [parseLine(l) for l in data.split('\n')]


def find_all_values(inputs: list[int], allow_concat: bool = False) -> list[int]:
    """Find all combinations of values."""
    if len(inputs) <= 1:
        return inputs
    *xs, x = inputs
    values = find_all_values(xs, allow_concat=allow_concat)
    results = []
    for v in values:
        results.append(v + x)
        results.append(v * x)
        if allow_concat:
            results.append(int(f'{v}{x}'))
    return results


def valid_equation(equation: Equation, allow_concat: bool = False) -> bool:
    """Determine if an equation is valid."""
    return equation.value in find_all_values(equation.inputs, allow_concat=allow_concat)


def main():
    with open(_DATA_PATH) as f:
        data = f.read()

    equations = parse(data)

    total = sum(
        eq.value for eq in equations
        if valid_equation(eq, allow_concat=False))
    print(f'1: {total}')

    total = sum(
        eq.value for eq in equations
        if valid_equation(eq, allow_concat=True))
    print(f'2: {total}')


if __name__ == '__main__':
    main()
