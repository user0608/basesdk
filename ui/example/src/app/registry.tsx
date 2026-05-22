import { defineRegistry } from "@basesdk/ui";
import { InventoryItemForm } from "./forms/InventoryItemForm";
import { SalesOrderForm } from "./forms/SalesOrderForm";
import { InventoryItemsPage } from "./views/InventoryItemsPage";
import { SalesOrdersPage } from "./views/SalesOrdersPage";

export const appRegistry = defineRegistry({
  pages: {
    "sales-orders-page": SalesOrdersPage,
    "inventory-items-page": InventoryItemsPage,
  },
  forms: {
    "sales-order-form": SalesOrderForm,
    "inventory-item-form": InventoryItemForm,
  },
});
