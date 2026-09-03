import htmx from "htmx.org/dist/htmx.esm";
import { Application } from "@hotwired/stimulus";
import { ProductRowController } from "./controllers/product_row_controller";
import { SearchBarController } from "./controllers/search_bar_controller";
import { DishFormController } from "./controllers/dish_form_controller";

window.htmx = htmx;
// htmx.logAll();

const stimulusApp = Application.start();
stimulusApp.debug = true;
window.C = stimulusApp;

stimulusApp.register("product-row", ProductRowController);
stimulusApp.register("search-bar", SearchBarController);
