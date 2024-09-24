ALTER TABLE games_categories
DROP CONSTRAINT games_categories_game_id_fkey,
ADD CONSTRAINT games_categories_game_id_fkey FOREIGN KEY (game_id) REFERENCES games(id) DEFERRABLE INITIALLY DEFERRED;