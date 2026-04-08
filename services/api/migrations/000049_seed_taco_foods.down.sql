-- Reverse seed: remove all TACO-sourced food items and related records.
-- Delete child tables first since there is no ON DELETE CASCADE.

DELETE FROM household_measures WHERE food_item_id IN (SELECT id FROM food_items WHERE source = 'taco');
DELETE FROM nutrition_facts WHERE food_item_id IN (SELECT id FROM food_items WHERE source = 'taco');
DELETE FROM food_items WHERE source = 'taco';
