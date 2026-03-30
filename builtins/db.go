package builtins

import (
	"encoding/json"
	"memmole/pkg/ast"
	"memmole/database"
	"memmole/pkg/object"
	"memmole/pkg/parser"
)

// EvalDBCall 处理 db.xxx(...) 调用
func EvalDBCall(exp *ast.MemberAccessExpression, args []object.Object, env *parser.Environment) object.Object {
	dbObjVal, exists := env.Get("db")
	if !exists {
		return &object.String{Value: "db 包未找到"}
	}
	dbo, ok := dbObjVal.(*database.DBObject)
	if !ok || dbo.Package == nil {
		return &object.String{Value: "无效的 db 对象"}
	}
	p := dbo.Package

	switch exp.Member.Value {
	case "connect":
		if len(args) != 2 {
			return &object.String{Value: "connect需要2个参数: backend, config_json"}
		}
		backend, ok1 := args[0].(*object.String)
		cfg, ok2 := args[1].(*object.String)
		if !ok1 || !ok2 {
			return &object.String{Value: "参数类型错误，应为 (string, string)"}
		}
		return p.Connect(backend.Value, cfg.Value)

	case "use":
		if len(args) != 1 {
			return &object.String{Value: "use需要1个参数: conn_id"}
		}
		id, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "use参数必须是整数"}
		}
		return p.Use(id.Value)

	case "create_table":
		if len(args) != 2 {
			return &object.String{Value: "create_table需要2个参数: name, schema_json"}
		}
		name, ok1 := args[0].(*object.String)
		schema, ok2 := args[1].(*object.String)
		if !ok1 || !ok2 {
			return &object.String{Value: "参数类型错误，应为 (string, string)"}
		}
		return p.CreateTable(name.Value, schema.Value)

	case "drop_table":
		if len(args) != 1 {
			return &object.String{Value: "drop_table需要1个参数: name"}
		}
		name, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "参数类型错误，应为 string"}
		}
		return p.DropTable(name.Value)

	case "insert":
		if len(args) != 2 {
			return &object.String{Value: "insert需要2个参数: table, row_json"}
		}
		tb, ok1 := args[0].(*object.String)
		data, ok2 := args[1].(*object.String)
		if !ok1 || !ok2 {
			return &object.String{Value: "参数类型错误，应为 (string, string)"}
		}
		return p.Insert(tb.Value, data.Value)

	case "update":
		if len(args) != 3 {
			return &object.String{Value: "update需要3个参数: table, where_json, data_json"}
		}
		tb, ok1 := args[0].(*object.String)
		wh, ok2 := args[1].(*object.String)
		data, ok3 := args[2].(*object.String)
		if !ok1 || !ok2 || !ok3 {
			return &object.String{Value: "参数类型错误，应为 (string, string, string)"}
		}
		return p.Update(tb.Value, wh.Value, data.Value)

	case "delete":
		if len(args) != 2 {
			return &object.String{Value: "delete需要2个参数: table, where_json"}
		}
		tb, ok1 := args[0].(*object.String)
		wh, ok2 := args[1].(*object.String)
		if !ok1 || !ok2 {
			return &object.String{Value: "参数类型错误，应为 (string, string)"}
		}
		return p.Delete(tb.Value, wh.Value)

	case "query":
		if len(args) < 2 || len(args) > 3 {
			return &object.String{Value: "query需要2~3个参数: table, where_json, [limit]"}
		}
		tb, ok1 := args[0].(*object.String)
		wh, ok2 := args[1].(*object.String)
		if !ok1 || !ok2 {
			return &object.String{Value: "前两个参数应为 (string, string)"}
		}
		var limit int64 = 0
		if len(args) == 3 {
			if li, ok := args[2].(*object.Integer); ok {
				limit = li.Value
			}
		}
		return p.Query(tb.Value, wh.Value, limit)

	case "query_raw":
		if len(args) < 1 {
			return &object.String{Value: "query_raw需要至少1个参数: sql"}
		}
		sql, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "query_raw第一个参数必须是字符串"}
		}
		return p.QueryRaw(sql.Value)

	case "execute_raw":
		if len(args) < 1 {
			return &object.String{Value: "execute_raw需要至少1个参数: sql"}
		}
		sql, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "execute_raw第一个参数必须是字符串"}
		}
		return p.ExecuteRaw(sql.Value)

	// 便捷 API
	case "create_table_from_struct":
		if len(args) != 1 {
			return &object.String{Value: "create_table_from_struct需要1个参数: struct_name"}
		}
		sname, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "参数类型错误，应为 string"}
		}
		if any, exists := env.Get(sname.Value); exists {
			if st, ok := any.(*object.StructType); ok {
				schema := make(map[string]string, len(st.Types))
				for k, v := range st.Types {
					schema[k] = v
				}
				b, _ := json.Marshal(schema)
				return p.CreateTable(sname.Value, string(b))
			}
		}
		return &object.String{Value: "未找到结构体类型或类型无效"}

	case "insert_kv":
		if len(args) < 3 || len(args)%2 == 0 {
			return &object.String{Value: "insert_kv参数: table, key, value, [key, value]..."}
		}
		tb, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "第一个参数必须是表名字符串"}
		}
		m := make(map[string]interface{})
		for i := 1; i < len(args); i += 2 {
			k, ok := args[i].(*object.String)
			if !ok {
				return &object.String{Value: "键必须是字符串"}
			}
			m[k.Value] = convertObjectToGo(args[i+1])
		}
		b, _ := json.Marshal(m)
		return p.Insert(tb.Value, string(b))

	case "query_kv":
		if len(args) < 1 {
			return &object.String{Value: "query_kv参数: table, [key, value]..."}
		}
		tb, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "第一个参数必须是表名字符串"}
		}
		where := make(map[string]interface{})
		if len(args) > 1 {
			if (len(args)-1)%2 != 0 {
				return &object.String{Value: "key/value 参数必须成对"}
			}
			for i := 1; i < len(args); i += 2 {
				k, ok := args[i].(*object.String)
				if !ok {
					return &object.String{Value: "键必须是字符串"}
				}
				where[k.Value] = convertObjectToGo(args[i+1])
			}
		}
		b, _ := json.Marshal(where)
		return p.Query(tb.Value, string(b), 0)

	case "update_by_id_kv":
		if len(args) < 3 || (len(args)-2)%2 != 0 {
			return &object.String{Value: "update_by_id_kv参数: table, id, key, value, ..."}
		}
		tb, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "第一个参数必须是表名字符串"}
		}
		id, ok := args[1].(*object.Integer)
		if !ok {
			return &object.String{Value: "第二个参数必须是id整数"}
		}
		data := make(map[string]interface{})
		for i := 2; i < len(args); i += 2 {
			k, ok := args[i].(*object.String)
			if !ok {
				return &object.String{Value: "键必须是字符串"}
			}
			data[k.Value] = convertObjectToGo(args[i+1])
		}
		where := map[string]interface{}{ "id": id.Value }
		wb, _ := json.Marshal(where)
		dbb, _ := json.Marshal(data)
		return p.Update(tb.Value, string(wb), string(dbb))

	case "delete_by_id":
		if len(args) != 2 {
			return &object.String{Value: "delete_by_id参数: table, id"}
		}
		tb, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "第一个参数必须是表名字符串"}
		}
		id, ok := args[1].(*object.Integer)
		if !ok {
			return &object.String{Value: "第二个参数必须是id整数"}
		}
		where := map[string]interface{}{ "id": id.Value }
		wb, _ := json.Marshal(where)
		return p.Delete(tb.Value, string(wb))

	case "insert_obj":
		if len(args) != 1 {
			return &object.String{Value: "insert_obj需要1个参数: struct_instance"}
		}
		if inst, ok := args[0].(*object.StructInstance); ok {
			m := map[string]interface{}{}
			for k, v := range inst.Fields {
				m[k] = convertObjectToGo(v)
			}
			b, _ := json.Marshal(m)
			return p.Insert(inst.TypeName, string(b))
		}
		return &object.String{Value: "参数必须是结构体实例"}

	case "update_obj_by_id":
		if len(args) != 1 {
			return &object.String{Value: "update_obj_by_id需要1个参数: struct_instance"}
		}
		if inst, ok := args[0].(*object.StructInstance); ok {
			where := map[string]interface{}{}
			data := map[string]interface{}{}
			for k, v := range inst.Fields {
				if k == "id" {
					if iv, ok := v.(*object.Integer); ok {
						where["id"] = iv.Value
					}
					continue
				}
				data[k] = convertObjectToGo(v)
			}
			if _, exists := where["id"]; !exists {
				return &object.String{Value: "对象未包含id，无法更新"}
			}
			wb, _ := json.Marshal(where)
			dbb, _ := json.Marshal(data)
			return p.Update(inst.TypeName, string(wb), string(dbb))
		}
		return &object.String{Value: "参数必须是结构体实例"}

	default:
		return &object.String{Value: "未知的db方法: " + exp.Member.Value}
	}
}

// convertObjectToGo 将解释器对象转为Go基础类型，便于JSON编码
func convertObjectToGo(o object.Object) interface{} {
	switch v := o.(type) {
	case *object.Integer:
		return v.Value
	case *object.Float:
		return v.Value
	case *object.String:
		return v.Value
	case *object.Boolean:
		return v.Value
	case *object.Null:
		return nil
	default:
		return v.Inspect()
	}
}
