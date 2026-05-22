import { defineMenuTree } from "@basesdk/ui";
import type { MenuTree, PageId } from "@basesdk/ui";
import { FiArchive, FiShoppingCart } from "react-icons/fi";
import { appRegistry } from "./registry";

type ExamplePageId = PageId<typeof appRegistry>;

export const appModules = defineMenuTree([
  {
    id: "sales",
    label: "Ventas",
    icon: FiShoppingCart,
    path: "/app/sales",
    order: 1,
    children: [
      {
        id: "sales-management",
        label: "Gestion",
        icon: FiShoppingCart,
        path: "management",
        order: 1,
        children: [
          {
            id: "sales-orders-list",
            label: "Ordenes",
            icon: FiShoppingCart,
            path: "orders",
            componentId: "sales-orders-page",
            order: 1,
          },
        ],
      },
    ],
  },
  {
    id: "inventory",
    label: "Inventarios",
    icon: FiArchive,
    path: "/app/inventory",
    order: 2,
    children: [
      {
        id: "inventory-catalog",
        label: "Catalogo",
        icon: FiArchive,
        path: "catalog",
        order: 1,
        children: [
          {
            id: "inventory-items-list",
            label: "Productos",
            icon: FiArchive,
            path: "items",
            componentId: "inventory-items-page",
            order: 1,
          },
        ],
      },
    ],
  },
] satisfies MenuTree<ExamplePageId, string>);
