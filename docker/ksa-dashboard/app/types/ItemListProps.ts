import { Item } from "./Item";

export type ItemListProps<T extends Item> = {
    data: Record<string, T[]>;
}