import ServicesClient from "./ServicesClient"
import { fetchWithAuth } from "@/lib/fetchWithAuth"
import { redirect } from "next/navigation"
import { ProductType } from "./types";

export default async function MitraPage(){
    try {
        const res = await fetchWithAuth("http://localhost:8080/api/products");
        if (!res.ok) redirect("/login");
    
        const productData: { products: ProductType[] } = await res.json();
    
        return <ServicesClient productData={productData.products} />;
      } catch (err) {
        redirect("/login");
      }

    return <ServicesClient productData={[]} />;
}