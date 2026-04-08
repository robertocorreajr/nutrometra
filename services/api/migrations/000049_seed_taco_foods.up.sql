-- Seed food catalog with data from the Brazilian TACO table
-- (Tabela Brasileira de Composição de Alimentos, 4th edition)
-- All nutritional values are per 100g of edible portion.
-- source = 'taco', tenant_id = NULL (global/system foods).

DO $$
DECLARE
  fid UUID;
BEGIN

  -- =====================================================================
  -- CEREAIS E DERIVADOS
  -- =====================================================================

  -- Arroz, integral, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Arroz, integral, cozido', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 124, 2.6, 25.8, 2.7, 1.0, 0.2, 1, 55, 5, 0.2, 43, 0.6, NULL);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 160);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'escumadeira', 90);

  -- Arroz, tipo 1, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Arroz, tipo 1, cozido', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 128, 2.5, 28.1, 1.6, 0.2, 0.1, 1, 32, 4, 0.1, 3, 0.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 160);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'escumadeira', 90);

  -- Aveia, flocos, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Aveia, flocos, crua', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 394, 13.9, 66.6, 9.1, 8.5, 1.5, 5, 336, 48, 4.4, 119, 2.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 15);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 80);

  -- Farinha de mandioca, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Farinha de mandioca, crua', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 361, 1.6, 87.9, 6.5, 0.3, 0.1, 1, 40, 1.0, 29);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 16);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 120);

  -- Farinha de trigo
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Farinha de trigo', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 360, 9.8, 75.1, 2.3, 1.4, 0.2, 1, 17, 1.0, 28, 0.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 13);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 120);

  -- Macarrão, trigo, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Macarrão, trigo, cozido', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 102, 3.4, 19.9, 1.5, 0.5, 0.1, 1, 6, 0.3, 14);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pegador cheio', 110);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'escumadeira', 75);

  -- Milho, verde, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Milho, verde, cozido', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg)
  VALUES (fid, 138, 6.6, 28.6, 3.9, 0.7, 0.1, 1, 282, 3, 0.5, 33, 0.5, 17);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'espiga média', 120);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 24);

  -- Pão, francês
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pão, francês', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 300, 8.0, 58.6, 2.3, 3.1, 0.8, 648, 22, 0.9, 21, 0.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 50);

  -- Pão, forma, integral
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pão, forma, integral', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 253, 9.4, 49.9, 6.9, 2.9, 0.6, 490, 159, 2.3, 42, 1.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 25);

  -- Pão de queijo
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pão de queijo', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg)
  VALUES (fid, 363, 5.1, 34.2, 0.5, 22.6, 7.1, 55, 427, 116, 0.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 20);

  -- Biscoito, cream cracker
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Biscoito, cream cracker', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 432, 9.2, 68.7, 2.1, 14.4, 3.7, 854, 16, 1.4, 18);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pacote individual (4 unidades)', 32);

  -- Cuscuz de milho cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Cuscuz de milho cozido', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 113, 2.4, 23.3, 1.7, 0.6, 0.1, 1, 7, 0.4, 8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 75);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);

  -- Tapioca (goma/fécula)
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Tapioca (goma/fécula)', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, calcium_mg, iron_mg)
  VALUES (fid, 343, 0.5, 84.4, 0.5, 0.1, 1, 12, 0.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 18);

  -- Granola
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Granola', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 421, 9.8, 71.5, 5.6, 12.5, 3.3, 36, 41, 3.0, 77, 1.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 12);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 60);

  -- Batata doce, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Batata doce, cozida', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 77, 0.6, 18.4, 2.2, 0.1, 6, 148, 17, 0.2, 11, 0.2, 11, 16.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 140);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);

  -- Mandioca, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Mandioca, cozida', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 125, 0.6, 30.1, 1.6, 0.3, 2, 100, 16, 0.2, 12, 11.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pedaço médio', 100);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 40);

  -- Batata inglesa, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Batata inglesa, cozida', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 52, 1.2, 11.9, 1.3, 0.0, 2, 188, 4, 0.3, 10, 6.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 140);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);

  -- Inhame, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Inhame, cozido', 'Cereais e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 97, 2.1, 23.2, 1.7, 0.1, 3, 204, 7, 0.2, 14);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pedaço médio', 100);

  -- =====================================================================
  -- VERDURAS, HORTALIÇAS E DERIVADOS
  -- =====================================================================

  -- Alface, crespa, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Alface, crespa, crua', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 11, 1.3, 1.7, 1.8, 0.2, 3, 267, 38, 0.4, 11, 234, 15.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'folha grande', 10);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'prato de sobremesa', 30);

  -- Rúcula, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Rúcula, crua', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 20, 2.6, 2.3, 1.6, 0.3, 4, 230, 117, 1.5, 20, 31.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'prato de sobremesa', 20);

  -- Espinafre, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Espinafre, cozido', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 22, 2.1, 2.6, 2.1, 0.2, 55, 311, 96, 0.5, 55, 376, 2.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);

  -- Couve, manteiga, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Couve, manteiga, crua', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 27, 2.9, 4.3, 3.1, 0.5, 6, 403, 131, 0.5, 27, 385, 96.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'folha grande', 30);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia (picada)', 18);

  -- Tomate, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Tomate, cru', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 15, 1.1, 3.1, 1.2, 0.1, 3, 222, 7, 0.2, 11, 54, 21.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 80);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 15);

  -- Cebola, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Cebola, crua', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 39, 1.7, 8.9, 1.7, 0.1, 1, 176, 15, 0.1, 9, 4.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 110);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia (picada)', 15);

  -- Alho, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Alho, cru', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 113, 7.0, 23.9, 4.3, 0.2, 6, 535, 14, 0.3, 21, 17.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'dente', 3);

  -- Cenoura, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Cenoura, crua', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 34, 1.3, 7.7, 3.2, 0.2, 3, 315, 23, 0.2, 11, 933, 5.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 65);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia (ralada)', 12);

  -- Beterraba, crua
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Beterraba, crua', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 49, 1.9, 11.1, 3.4, 0.1, 14, 375, 18, 0.3, 18, 3.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 100);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia (ralada)', 15);

  -- Brócolis, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Brócolis, cozido', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 25, 2.1, 4.4, 3.4, 0.3, 3, 178, 51, 0.5, 14, 111, 42.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'ramo médio', 15);

  -- Couve-flor, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Couve-flor, cozida', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 19, 1.2, 3.9, 2.1, 0.2, 1, 176, 16, 0.3, 8, 23.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'ramo médio', 20);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 35);

  -- Pepino, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pepino, cru', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 10, 0.9, 2.0, 1.1, 0.0, 1, 154, 11, 0.2, 8, 5.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 130);

  -- Abobrinha, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Abobrinha, cozida', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 15, 0.8, 3.3, 1.2, 0.1, 1, 183, 13, 0.2, 10, 6.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);

  -- Berinjela, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Berinjela, cozida', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 19, 0.7, 4.5, 2.5, 0.1, 1, 94, 8, 0.2, 9);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);

  -- Pimentão verde, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pimentão verde, cru', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 21, 1.0, 4.9, 2.6, 0.1, 1, 152, 6, 0.2, 8, 16, 100.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 70);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 7);

  -- Pimentão vermelho, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pimentão vermelho, cru', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 27, 1.2, 5.8, 1.4, 0.3, 1, 190, 6, 0.3, 9, 122, 158.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 70);

  -- Repolho, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Repolho, cru', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 17, 0.9, 3.9, 1.9, 0.1, 6, 170, 35, 0.4, 12, 32.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia (picado)', 10);

  -- Quiabo, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Quiabo, cozido', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 22, 1.4, 4.4, 3.0, 0.2, 2, 113, 62, 0.3, 33, 9.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);

  -- Abóbora, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Abóbora, cozida', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 28, 0.8, 7.0, 1.6, 0.1, 1, 119, 15, 0.2, 6, 188, 7.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 36);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 70);

  -- Chuchu, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Chuchu, cozido', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 17, 0.4, 4.0, 1.5, 0.1, 1, 96, 10, 0.1, 8, 9.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);

  -- Vagem, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Vagem, cozida', 'Verduras, hortaliças e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 25, 1.5, 5.2, 2.4, 0.1, 2, 129, 28, 0.5, 14, 4.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);

  -- =====================================================================
  -- FRUTAS E DERIVADOS
  -- =====================================================================

  -- Banana, prata
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Banana, prata', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 98, 1.3, 26.0, 2.0, 0.1, 1, 358, 8, 0.4, 26, 21.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 86);

  -- Banana, nanica
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Banana, nanica', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 92, 1.4, 23.8, 1.9, 0.1, 1, 376, 3, 0.3, 28, 5.9);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 101);

  -- Maçã, fuji
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Maçã, fuji', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 56, 0.3, 15.2, 1.3, 0.0, 1, 75, 2, 0.1, 2, 2.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 130);

  -- Laranja, pera
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Laranja, pera', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 37, 1.0, 8.9, 0.8, 0.1, 1, 163, 22, 0.1, 11, 53.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 137);

  -- Mamão, papaia
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Mamão, papaia', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 40, 0.5, 10.4, 1.0, 0.1, 3, 222, 25, 0.2, 17, 37, 82.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 170);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 350);

  -- Manga, tommy
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Manga, tommy', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 51, 0.4, 12.8, 1.6, 0.2, 1, 148, 9, 0.1, 7, 210, 17.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 220);

  -- Abacaxi, pérola
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Abacaxi, pérola', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 48, 0.9, 12.3, 1.0, 0.1, 1, 131, 22, 0.3, 18, 34.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 75);

  -- Morango, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Morango, cru', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 30, 0.9, 6.8, 1.7, 0.3, 1, 184, 11, 0.3, 10, 63.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 12);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 140);

  -- Melancia
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Melancia', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 33, 0.9, 8.1, 0.1, 0.0, 1, 104, 8, 0.2, 7, 23, 6.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 200);

  -- Melão
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Melão', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 29, 0.7, 7.5, 0.3, 0.0, 11, 216, 3, 0.2, 6, 1.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 115);

  -- Uva, itália
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Uva, itália', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 53, 0.7, 13.6, 0.9, 0.2, 1, 162, 7, 0.1, 4, 2.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'cacho médio (15 bagos)', 100);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'bago', 7);

  -- Abacate
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Abacate', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 96, 1.2, 6.0, 6.3, 8.4, 2.0, 1, 206, 8, 0.2, 15, 8.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 30);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média (polpa)', 170);

  -- Limão, tahiti
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Limão, tahiti', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 32, 0.9, 11.1, 1.2, 0.1, 1, 128, 51, 0.2, 5, 38.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 65);

  -- Goiaba, vermelha
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Goiaba, vermelha', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 54, 1.1, 13.0, 6.2, 0.4, 1, 198, 4, 0.2, 7, 69, 80.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 170);

  -- Pera
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pera', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 53, 0.6, 14.0, 3.0, 0.1, 1, 116, 9, 0.1, 5, 3.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 133);

  -- Kiwi
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Kiwi', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 51, 1.3, 11.5, 2.7, 0.6, 1, 269, 24, 0.3, 15, 70.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 76);

  -- Acerola
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Acerola', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 33, 0.9, 8.0, 1.5, 0.2, 3, 165, 13, 0.2, 13, 40, 941.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 120);

  -- Maracujá, polpa
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Maracujá, polpa', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 68, 2.0, 12.3, 1.1, 2.1, 4, 338, 5, 0.6, 28, 19.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média (polpa)', 55);

  -- Coco, polpa, cru
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Coco, polpa, cru', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 406, 3.7, 10.4, 5.4, 40.2, 33.7, 5, 268, 6, 1.5, 47, 0.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pedaço médio', 40);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia (ralado)', 8);

  -- Tangerina
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Tangerina', 'Frutas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_a_mcg, vitamin_c_mg)
  VALUES (fid, 38, 0.8, 9.6, 0.8, 0.1, 1, 176, 12, 0.1, 11, 41, 48.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 135);

  -- =====================================================================
  -- CARNES E DERIVADOS
  -- =====================================================================

  -- Frango, peito, sem pele, grelhado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Frango, peito, sem pele, grelhado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 159, 32.0, 0, 0, 2.5, 0.8, 89, 46, 340, 4, 0.3, 29, 1.0, 0.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'filé médio', 120);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'filé pequeno', 80);

  -- Frango, coxa, sem pele, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Frango, coxa, sem pele, cozida', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 163, 26.2, 0, 0, 5.8, 1.6, 154, 59, 231, 8, 0.6, 23, 2.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 65);

  -- Frango, sobrecoxa, sem pele, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Frango, sobrecoxa, sem pele, cozida', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 185, 25.6, 0, 0, 8.5, 2.4, 118, 68, 204, 10, 0.7, 20, 2.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade média', 100);

  -- Boi, acém, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Boi, acém, cozido', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 212, 26.7, 0, 0, 11.2, 4.5, 81, 44, 217, 6, 2.6, 17, 7.3, 2.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pedaço médio', 60);

  -- Boi, alcatra, grelhada
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Boi, alcatra, grelhada', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 170, 32.4, 0, 0, 3.9, 1.5, 72, 54, 356, 3, 2.9, 25, 4.4, 2.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'bife médio', 100);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'bife pequeno', 65);

  -- Boi, patinho, grelhado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Boi, patinho, grelhado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 166, 32.8, 0, 0, 3.2, 1.3, 79, 45, 348, 3, 3.0, 24, 5.5, 1.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'bife médio', 100);

  -- Boi, filé mignon, grelhado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Boi, filé mignon, grelhado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 165, 32.8, 0, 0, 3.2, 1.2, 92, 51, 365, 3, 3.2, 27, 4.5, 2.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'bife médio', 120);

  -- Boi, coxão mole, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Boi, coxão mole, cozido', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 199, 32.4, 0, 0, 7.3, 2.8, 86, 47, 280, 4, 2.7, 22, 6.0, 2.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pedaço médio', 75);

  -- Carne moída, refogada
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Carne moída, refogada', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 212, 26.0, 0, 0, 11.5, 4.6, 82, 40, 235, 5, 2.4, 18, 5.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'escumadeira', 65);

  -- Boi, fígado, grelhado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Boi, fígado, grelhado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg, vitamin_b12_mcg)
  VALUES (fid, 225, 29.5, 4.3, 0, 9.5, 3.2, 397, 72, 354, 6, 5.8, 19, 4.3, 4968, 70.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'bife médio', 100);

  -- Porco, lombo, assado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Porco, lombo, assado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_b12_mcg)
  VALUES (fid, 210, 30.2, 0, 0, 9.5, 3.4, 78, 54, 356, 5, 0.8, 22, 2.3, 0.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 50);

  -- Porco, costela, assada
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Porco, costela, assada', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 292, 23.2, 0, 0, 21.8, 7.8, 85, 60, 264, 14, 0.8, 17, 3.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pedaço médio', 80);

  -- Peru, peito, sem pele, assado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Peru, peito, sem pele, assado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 128, 25.5, 0, 0, 2.3, 0.7, 100, 60, 291, 11, 0.5, 26, 1.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 35);

  -- Linguiça, de porco, grelhada
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Linguiça, de porco, grelhada', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 296, 18.5, 0.7, 0, 24.3, 9.1, 85, 1085, 8, 0.9, 12, 2.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'gomo médio', 50);

  -- Presunto, magro
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Presunto, magro', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 108, 17.4, 1.8, 0, 3.5, 1.3, 48, 1228, 6, 0.5, 17, 1.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 15);

  -- Peito de peru defumado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Peito de peru defumado', 'Carnes e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 118, 22.2, 1.0, 0, 2.5, 0.9, 52, 1121, 8, 0.3, 24, 1.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 15);

  -- =====================================================================
  -- LEITE E DERIVADOS
  -- =====================================================================

  -- Leite de vaca, integral
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Leite de vaca, integral', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg)
  VALUES (fid, 61, 3.2, 4.7, 0, 3.3, 1.9, 12, 61, 164, 123, 0.1, 10, 0.4, 46);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'copo (200ml)', 200);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 15);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 240);

  -- Leite de vaca, desnatado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Leite de vaca, desnatado', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg)
  VALUES (fid, 35, 3.4, 4.9, 0, 0.2, 0.1, 3, 53, 166, 134, 0.1, 12, 0.4, 3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'copo (200ml)', 200);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 240);

  -- Iogurte, natural
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Iogurte, natural', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 51, 4.1, 5.2, 0, 1.6, 1.0, 6, 52, 187, 143, 0.1, 12, 0.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pote (170g)', 170);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 20);

  -- Iogurte, natural, desnatado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Iogurte, natural, desnatado', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 40, 4.1, 5.8, 0, 0.2, 0.1, 1, 56, 193, 149, 0.1, 13);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'pote (170g)', 170);

  -- Queijo, mussarela
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Queijo, mussarela', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg)
  VALUES (fid, 330, 22.6, 3.0, 0, 25.2, 14.8, 86, 579, 62, 570, 0.2, 21, 2.8, 192);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 15);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia grossa', 30);

  -- Queijo, prato
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Queijo, prato', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg)
  VALUES (fid, 360, 22.7, 1.9, 0, 29.1, 17.1, 93, 585, 82, 760, 0.3, 25, 3.0, 211);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia', 20);

  -- Queijo, cottage
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Queijo, cottage', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 98, 13.7, 3.0, 0, 3.4, 2.1, 13, 375, 60, 0.1, 5, 0.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 20);

  -- Queijo, ricota
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Queijo, ricota', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 140, 12.6, 3.7, 0, 8.1, 5.1, 49, 204, 253, 0.2, 10, 0.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 30);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 20);

  -- Queijo, minas frescal
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Queijo, minas frescal', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 264, 17.4, 3.2, 0, 20.2, 12.4, 62, 440, 579, 0.2, 14, 1.9);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'fatia média', 30);

  -- Requeijão, cremoso
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Requeijão, cremoso', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 257, 7.6, 2.5, 0, 24.0, 14.9, 70, 353, 170, 0.2, 8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 15);

  -- Manteiga
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Manteiga', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg, vitamin_a_mcg)
  VALUES (fid, 726, 0.4, 0.0, 0, 82.4, 51.2, 201, 579, 12, 0.1, 685);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa', 10);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'ponta de faca', 3);

  -- Creme de leite
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Creme de leite', 'Leite e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, calcium_mg, iron_mg)
  VALUES (fid, 209, 2.0, 3.6, 0, 21.0, 13.1, 76, 34, 65, 0.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 15);

  -- =====================================================================
  -- LEGUMINOSAS E DERIVADOS
  -- =====================================================================

  -- Feijão, carioca, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Feijão, carioca, cozido', 'Leguminosas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 76, 4.8, 13.6, 8.5, 0.5, 0.1, 2, 256, 27, 1.3, 40, 0.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'concha média', 80);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 17);

  -- Feijão, preto, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Feijão, preto, cozido', 'Leguminosas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 77, 4.5, 14.0, 8.4, 0.5, 0.1, 2, 294, 29, 1.5, 50, 0.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'concha média', 80);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 17);

  -- Lentilha, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Lentilha, cozida', 'Leguminosas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 93, 6.3, 16.3, 7.9, 0.5, 0.1, 2, 220, 16, 1.5, 22, 0.9);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'concha média', 80);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 17);

  -- Grão de bico, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Grão de bico, cozido', 'Leguminosas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 130, 6.7, 20.0, 5.1, 2.6, 0.3, 6, 195, 46, 2.1, 36, 1.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'concha média', 80);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 17);

  -- Soja, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Soja, cozida', 'Leguminosas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 151, 14.0, 7.1, 5.6, 7.6, 1.1, 1, 449, 83, 2.5, 62, 1.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'concha média', 80);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 17);

  -- Ervilha, seca, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Ervilha, seca, cozida', 'Leguminosas e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 88, 6.1, 15.7, 5.1, 0.4, 0.1, 3, 204, 16, 1.1, 22, 0.8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 17);

  -- =====================================================================
  -- OVOS E DERIVADOS
  -- =====================================================================

  -- Ovo, de galinha, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Ovo, de galinha, cozido', 'Ovos e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg, vitamin_b12_mcg)
  VALUES (fid, 146, 13.3, 0.6, 0, 9.5, 3.1, 397, 137, 136, 49, 1.5, 11, 1.1, 155, 1.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 50);

  -- Ovo, frito
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Ovo, frito', 'Ovos e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg)
  VALUES (fid, 240, 15.6, 0.6, 0, 19.7, 5.8, 516, 239, 133, 52, 1.7, 12, 1.4, 165);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 46);

  -- Clara de ovo, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Clara de ovo, cozida', 'Ovos e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 48, 10.5, 0.8, 0, 0.0, 0, 171, 153, 6, 0.1, 11, 0.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 33);

  -- Gema de ovo, cozida
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Gema de ovo, cozida', 'Ovos e derivados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg, vitamin_a_mcg, vitamin_b12_mcg)
  VALUES (fid, 352, 15.9, 2.5, 0, 30.8, 10.0, 1272, 51, 113, 120, 3.1, 8, 2.6, 340, 2.7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 17);

  -- =====================================================================
  -- PESCADOS E FRUTOS DO MAR
  -- =====================================================================

  -- Salmão, grelhado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Salmão, grelhado', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 243, 23.8, 0, 0, 15.6, 3.4, 57, 87, 269, 10, 0.3, 25, 0.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'posta média', 120);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'filé médio', 100);

  -- Atum, em conserva
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Atum, em conserva', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 166, 26.2, 0, 0, 6.4, 1.5, 41, 396, 187, 8, 0.9, 27, 0.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'lata (170g escorrido ~120g)', 120);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 25);

  -- Tilápia, filé, grelhado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Tilápia, filé, grelhado', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 128, 26.0, 0, 0, 2.7, 0.9, 50, 45, 302, 7, 0.4, 26, 0.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'filé médio', 120);

  -- Sardinha, em conserva
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Sardinha, em conserva', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 214, 23.9, 0, 0, 12.7, 3.7, 61, 499, 270, 550, 1.8, 30, 1.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 20);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'lata (125g escorrido ~85g)', 85);

  -- Bacalhau, salgado, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Bacalhau, salgado, cozido', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 138, 29.0, 0, 0, 2.1, 0.4, 68, 254, 302, 13, 0.4, 30, 0.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'posta média', 100);

  -- Camarão, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Camarão, cozido', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 90, 18.4, 0, 0, 1.5, 0.4, 129, 365, 175, 47, 0.8, 30, 0.9);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade grande', 15);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 80);

  -- Merluza, filé, cozido
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Merluza, filé, cozido', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 122, 18.9, 0, 0, 5.0, 1.2, 65, 84, 241, 14, 0.2, 18, 0.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'filé médio', 100);

  -- Pescada, filé, frito
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Pescada, filé, frito', 'Pescados e frutos do mar', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 183, 19.9, 5.0, 0, 9.0, 2.2, 56, 187, 243, 18, 0.4, 22, 0.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'filé médio', 100);

  -- =====================================================================
  -- GORDURAS E ÓLEOS
  -- =====================================================================

  -- Azeite de oliva
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Azeite de oliva', 'Gorduras e óleos', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg)
  VALUES (fid, 884, 0, 0, 0, 100, 14.0, 0, 1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa', 13);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 5);

  -- Óleo de soja
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Óleo de soja', 'Gorduras e óleos', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg)
  VALUES (fid, 884, 0, 0, 0, 100, 15.2, 0, 0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa', 13);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 5);

  -- Óleo de coco
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Óleo de coco', 'Gorduras e óleos', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg)
  VALUES (fid, 862, 0, 0, 0, 100, 82.0, 0, 0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa', 13);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 5);

  -- Margarina
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Margarina', 'Gorduras e óleos', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, trans_fat_g, cholesterol_mg, sodium_mg, vitamin_a_mcg)
  VALUES (fid, 720, 0.1, 0.1, 0, 80.0, 17.4, 6.2, 0, 710, 610);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'ponta de faca', 3);

  -- =====================================================================
  -- NOZES E SEMENTES
  -- =====================================================================

  -- Castanha do Pará
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Castanha do Pará', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 643, 14.5, 12.3, 7.9, 63.5, 15.1, 1, 600, 146, 2.3, 376, 4.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 4);

  -- Castanha de caju
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Castanha de caju', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 570, 18.5, 29.1, 3.7, 46.3, 8.1, 7, 565, 29, 5.4, 260, 5.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'punhado (10 unidades)', 30);

  -- Amendoim, torrado
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Amendoim, torrado', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 606, 27.2, 12.5, 7.8, 49.4, 7.5, 5, 580, 49, 1.8, 150, 3.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 15);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'punhado', 30);

  -- Semente de chia
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Semente de chia', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 486, 16.5, 42.1, 34.4, 30.7, 3.3, 16, 407, 631, 7.7, 335, 4.6);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 12);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 4);

  -- Semente de linhaça
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Semente de linhaça', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 495, 14.1, 43.3, 33.5, 32.3, 3.0, 30, 813, 211, 4.7, 362, 4.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 12);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 4);

  -- Semente de girassol
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Semente de girassol', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 570, 20.8, 20.0, 8.6, 51.5, 5.2, 9, 645, 70, 5.3, 325, 5.0);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa cheia', 12);

  -- Amêndoa
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Amêndoa', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 581, 21.2, 21.7, 12.5, 49.9, 3.8, 1, 705, 248, 3.4, 268, 3.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 1.2);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'punhado (10 unidades)', 12);

  -- Noz
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Noz', 'Nozes e sementes', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 620, 14.4, 18.4, 7.5, 59.0, 5.6, 2, 441, 84, 2.6, 158, 3.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'unidade', 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'punhado (5 unidades)', 25);

  -- =====================================================================
  -- PRODUTOS AÇUCARADOS
  -- =====================================================================

  -- Açúcar, cristal
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Açúcar, cristal', 'Produtos açucarados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, calcium_mg, iron_mg)
  VALUES (fid, 387, 0, 99.6, 0, 0, 1, 3, 0.1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa', 15);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara', 160);

  -- Mel de abelha
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Mel de abelha', 'Produtos açucarados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 309, 0.3, 84.0, 0, 0, 5, 54, 4, 0.3, 1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de chá', 7);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'colher de sopa', 21);

  -- Chocolate ao leite
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Chocolate ao leite', 'Produtos açucarados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, cholesterol_mg, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 540, 7.2, 59.4, 2.9, 30.3, 17.4, 16, 61, 353, 169, 1.9, 58, 1.4);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'barra pequena (25g)', 25);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'quadrado', 5);

  -- Chocolate amargo (70%)
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Chocolate amargo (70%)', 'Produtos açucarados', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, saturated_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, zinc_mg)
  VALUES (fid, 530, 10.9, 38.3, 10.0, 38.3, 22.8, 8, 567, 56, 8.0, 176, 2.5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'barra pequena (25g)', 25);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'quadrado', 5);

  -- =====================================================================
  -- BEBIDAS
  -- =====================================================================

  -- Café, infusão
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Café, infusão', 'Bebidas', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 3, 0.4, 0.3, 0, 0, 3, 96, 3, 0.1, 5);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara (50ml)', 50);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara de chá (200ml)', 200);

  -- Suco de laranja, natural
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Suco de laranja, natural', 'Bebidas', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg, vitamin_c_mg)
  VALUES (fid, 45, 0.6, 10.6, 0.2, 0.1, 1, 186, 10, 0.1, 10, 73.3);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'copo (200ml)', 200);

  -- Água de coco
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Água de coco', 'Bebidas', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 22, 0, 5.3, 0, 0, 2, 162, 22, 0.1, 8);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'copo (200ml)', 200);

  -- Chá verde, infusão
  INSERT INTO food_items (name, food_group, serving_size_g, serving_label, source, tenant_id)
  VALUES ('Chá verde, infusão', 'Bebidas', 100, '100g', 'taco', NULL)
  RETURNING id INTO fid;
  INSERT INTO nutrition_facts (food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, total_fat_g, sodium_mg, potassium_mg, calcium_mg, iron_mg, magnesium_mg)
  VALUES (fid, 1, 0.2, 0, 0, 0, 1, 8, 1, 0.0, 1);
  INSERT INTO household_measures (food_item_id, label, grams) VALUES (fid, 'xícara de chá (200ml)', 200);

END $$;
