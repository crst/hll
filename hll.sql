WITH random_numbers AS (
  SELECT (random() * (1 << 16)) :: INT AS n FROM generate_series(1, 1 << 14)
),
params AS (
  SELECT
    b,
    2^b :: INT AS m,
    CASE
      WHEN b = 4 THEN 0.673
      WHEN b = 5 THEN 0.697
      WHEN b = 6 THEN 0.709
      ELSE 0.7213 / (1.0 + (1.079 / 2^b))
    END AS alpha,
    1.04 / sqrt(2^b) AS expected_accuracy
  FROM (SELECT 8 AS b) p
),
registers AS (
  SELECT
    (x & ((1 << b) - 1)) AS register,
    max((62 - b) - floor(log(2, ((x & ((1 << 62) - 1)) >> b)))) :: INT AS p
  FROM (
    SELECT
      ('x' || substring(md5(n :: VARCHAR), 1, 16)) :: BIT(64) :: BIGINT AS x
    FROM random_numbers
  ) t
  CROSS JOIN params
  GROUP BY register
),
approximation AS (
  SELECT
    ((alpha * m * m) / Z) :: BIGINT AS num_estimated_distinct_elements
  FROM (
    SELECT
      sum(1.0 / (1 << p)) AS Z
    FROM registers
  ) t
  CROSS JOIN params
)
SELECT
  num_distinct_elements,
  num_estimated_distinct_elements,
  (100 * ((num_estimated_distinct_elements - num_distinct_elements) :: FLOAT
    / num_estimated_distinct_elements)) :: NUMERIC(5, 2) AS p_error,
  (100 * expected_accuracy) :: NUMERIC(5, 2) AS p_expected_accuracy
FROM (SELECT count(DISTINCT n) AS num_distinct_elements FROM random_numbers) exact_result
CROSS JOIN approximation approximation
CROSS JOIN params;
