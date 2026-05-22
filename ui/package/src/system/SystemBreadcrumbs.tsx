import { Link } from "react-router-dom";
import { FiChevronRight } from "react-icons/fi";

export type SystemBreadcrumbItem = {
  label: string;
  to?: string;
};

export const SystemBreadcrumbs = ({ items, className = "" }: { items: SystemBreadcrumbItem[]; className?: string }) => {
  return (
    <nav aria-label="Breadcrumb" className={["flex min-w-0 items-center gap-1 text-xs text-ui-text-soft", className].join(" ").trim()}>
      {items.map((item, index) => {
        const isLast = index === items.length - 1;

        return (
          <span key={`${item.label}-${index}`} className="flex min-w-0 items-center gap-1">
            {item.to && !isLast ? (
              <Link to={item.to} className="truncate font-medium transition-colors hover:text-ui-primary">
                {item.label}
              </Link>
            ) : (
              <span className={isLast ? "truncate font-semibold text-ui-text-muted" : "truncate"}>{item.label}</span>
            )}
            {!isLast && <FiChevronRight size={13} className="shrink-0" />}
          </span>
        );
      })}
    </nav>
  );
};
