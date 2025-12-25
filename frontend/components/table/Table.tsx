"use client";
import styles from './Table.module.css';

interface Header {
  key: string;
  label: string;
  align?: 'left' | 'right' | 'center';
}

interface TableProps {
  headers: Header[];
  data: any[];
  onUpdate?: (id: any) => void;
  onDelete?: (id: any) => void;
  onDetail?: (id: any) => void;
  renderCustomCell?: (key: string, value: any, item: any) => React.ReactNode;
}

export default function Table({ 
  headers, 
  data, 
  onUpdate, 
  onDelete, 
  onDetail,
  renderCustomCell 
}: TableProps) {
  return (
    <div className={styles.tableContainer}>
      <table className={styles.table}>
        <thead>
          <tr>
            {headers.map((header) => (
              <th key={header.key} style={{ textAlign: header.align || 'left' }}>
                {header.label}
              </th>
            ))}
            {(onUpdate || onDelete || onDetail) && (
              <th style={{ textAlign: 'right' }}>Actions</th>
            )}
          </tr>
        </thead>
        <tbody>
          {data.map((item, index) => (
            <tr key={item.id || index}>
              {headers.map((header) => (
                <td key={header.key} style={{ textAlign: header.align || 'left' }}>
                  {renderCustomCell 
                    ? renderCustomCell(header.key, item[header.key], item) 
                    : item[header.key]}
                </td>
              ))}
              {(onUpdate || onDelete || onDetail) && (
                <td style={{ textAlign: 'right' }}>
                  <div className={styles.actionGroup}>
                    {onDetail && (
                      <button className={styles.btnDetail} onClick={() => onDetail(item.id)}>View</button>
                    )}
                    {onUpdate && (
                      <button className={styles.btnUpdate} onClick={() => onUpdate(item.id)}>Update</button>
                    )}
                    {onDelete && (
                      <button className={styles.btnDelete} onClick={() => onDelete(item.id)}>Delete</button>
                    )}
                  </div>
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}