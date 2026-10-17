<template>
  <div class="system-dict-container">
    <el-row :gutter="15">
      <!-- 左侧: 字典类型 -->
      <el-col :span="10">
        <el-card shadow="never">
          <div class="table-toolbar">
            <el-button
              v-if="checkPermission(['system:dict:create'])"
              type="primary"
              size="small"
              @click="handleTypeCreate"
            >
              新增类型
            </el-button>
          </div>
          <el-table
            v-loading="typeLoading"
            :data="typeList"
            border
            highlight-current-row
            @current-change="handleTypeSelect"
          >
            <el-table-column prop="name" label="字典名称" min-width="100" />
            <el-table-column prop="dictType" label="类型键" min-width="120">
              <template #default="{ row }">
                <el-tag size="small">{{ row.dictType }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="70" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                  {{ row.status === 1 ? "启用" : "停用" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="130" align="center">
              <template #default="{ row }">
                <el-button
                  v-if="checkPermission(['system:dict:update'])"
                  type="primary"
                  link
                  @click.stop="handleTypeEdit(row)"
                >
                  编辑
                </el-button>
                <el-button
                  v-if="checkPermission(['system:dict:delete'])"
                  type="danger"
                  link
                  @click.stop="handleTypeDelete(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 右侧: 字典数据 (跟随左侧选中类型) -->
      <el-col :span="14">
        <el-card shadow="never">
          <div class="table-toolbar">
            <el-button
              v-if="checkPermission(['system:dict:create']) && currentType"
              type="primary"
              size="small"
              @click="handleDataCreate"
            >
              新增数据
            </el-button>
            <span v-if="currentType" class="current-type">
              当前类型: {{ currentType.name }} ({{ currentType.dictType }})
            </span>
          </div>
          <el-table v-loading="dataLoading" :data="dataList" border>
            <el-table-column prop="label" label="标签" min-width="100">
              <template #default="{ row }">
                <el-tag v-if="row.tagClass" :type="row.tagClass" size="small">
                  {{ row.label }}
                </el-tag>
                <span v-else>{{ row.label }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="value" label="键值" min-width="90" />
            <el-table-column prop="sort" label="排序" width="70" align="center" />
            <el-table-column label="状态" width="70" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                  {{ row.status === 1 ? "启用" : "停用" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="100" />
            <el-table-column label="操作" width="130" align="center">
              <template #default="{ row }">
                <el-button
                  v-if="checkPermission(['system:dict:update'])"
                  type="primary"
                  link
                  @click="handleDataEdit(row)"
                >
                  编辑
                </el-button>
                <el-button
                  v-if="checkPermission(['system:dict:delete'])"
                  type="danger"
                  link
                  @click="handleDataDelete(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 字典类型对话框 -->
    <el-dialog
      v-model="typeDialogVisible"
      :title="typeForm.id ? '编辑字典类型' : '新增字典类型'"
      width="480px"
      destroy-on-close
    >
      <el-form ref="typeFormRef" :model="typeForm" :rules="typeRules" label-width="90px">
        <el-form-item label="字典名称" prop="name">
          <el-input v-model="typeForm.name" placeholder="请输入字典名称" />
        </el-form-item>
        <el-form-item label="类型键" prop="dictType">
          <el-input v-model="typeForm.dictType" placeholder="请输入类型键(唯一, 如 sys_status)" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="typeForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="typeForm.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="typeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleTypeSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 字典数据对话框 -->
    <el-dialog
      v-model="dataDialogVisible"
      :title="dataForm.id ? '编辑字典数据' : '新增字典数据'"
      width="480px"
      destroy-on-close
    >
      <el-form ref="dataFormRef" :model="dataForm" :rules="dataRules" label-width="90px">
        <el-form-item label="标签" prop="label">
          <el-input v-model="dataForm.label" placeholder="请输入标签(展示用)" />
        </el-form-item>
        <el-form-item label="键值" prop="value">
          <el-input v-model="dataForm.value" placeholder="请输入键值(存储用)" />
        </el-form-item>
        <el-form-item label="标签样式">
          <el-select v-model="dataForm.tagClass" placeholder="默认(无样式)" clearable style="width: 100%">
            <el-option label="primary" value="primary" />
            <el-option label="success" value="success" />
            <el-option label="warning" value="warning" />
            <el-option label="danger" value="danger" />
            <el-option label="info" value="info" />
          </el-select>
        </el-form-item>
        <el-form-item label="显示顺序">
          <el-input-number v-model="dataForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="dataForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="dataForm.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dataDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleDataSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
  import { ElMessage, ElMessageBox } from "element-plus";
  import checkPermission from "@/utils/permission";
  import {
    getDictTypeList,
    createDictType,
    updateDictType,
    deleteDictType,
    getDictDataList,
    createDictData,
    updateDictData,
    deleteDictData,
  } from "@/api/system";

  export default {
    name: "SystemDict",
    data() {
      return {
        typeLoading: false,
        dataLoading: false,
        submitting: false,
        typeList: [],
        dataList: [],
        currentType: null,
        typeDialogVisible: false,
        typeForm: this.getDefaultTypeForm(),
        typeRules: {
          name: [{ required: true, message: "请输入字典名称", trigger: "blur" }],
          dictType: [{ required: true, message: "请输入类型键", trigger: "blur" }],
        },
        dataDialogVisible: false,
        dataForm: this.getDefaultDataForm(),
        dataRules: {
          label: [{ required: true, message: "请输入标签", trigger: "blur" }],
          value: [{ required: true, message: "请输入键值", trigger: "blur" }],
        },
      };
    },
    created() {
      this.fetchTypes();
    },
    methods: {
      checkPermission,
      getDefaultTypeForm() {
        return { id: undefined, name: "", dictType: "", status: 1, remark: "" };
      },
      getDefaultDataForm() {
        return {
          id: undefined,
          dictType: "",
          label: "",
          value: "",
          tagClass: "",
          sort: 1,
          status: 1,
          remark: "",
        };
      },
      // 查询字典类型列表
      async fetchTypes() {
        this.typeLoading = true;
        try {
          const { data } = await getDictTypeList();
          this.typeList = data.list || [];
        } finally {
          this.typeLoading = false;
        }
      },
      // 查询当前类型下的字典数据
      async fetchData() {
        if (!this.currentType) return;
        this.dataLoading = true;
        try {
          const { data } = await getDictDataList({ dictType: this.currentType.dictType });
          this.dataList = data.list || [];
        } finally {
          this.dataLoading = false;
        }
      },
      // 选中类型联动右侧数据
      handleTypeSelect(row) {
        this.currentType = row;
        if (row) this.fetchData();
        else this.dataList = [];
      },
      handleTypeCreate() {
        this.typeForm = this.getDefaultTypeForm();
        this.typeDialogVisible = true;
      },
      handleTypeEdit(row) {
        this.typeForm = {
          id: row.id,
          name: row.name,
          dictType: row.dictType,
          status: row.status,
          remark: row.remark,
        };
        this.typeDialogVisible = true;
      },
      handleTypeSubmit() {
        this.$refs.typeFormRef.validate(async (valid) => {
          if (!valid) return;
          this.submitting = true;
          try {
            if (this.typeForm.id) {
              await updateDictType(this.typeForm);
              ElMessage.success("修改成功");
            } else {
              await createDictType(this.typeForm);
              ElMessage.success("新增成功");
            }
            this.typeDialogVisible = false;
            this.fetchTypes();
          } finally {
            this.submitting = false;
          }
        });
      },
      // 删除字典类型 (后端级联删除其下数据)
      handleTypeDelete(row) {
        ElMessageBox.confirm(
          `确认删除字典类型 [${row.name}] 吗? 其下字典数据将一并删除`,
          "提示",
          { type: "warning" }
        ).then(async () => {
          await deleteDictType({ ids: [row.id] });
          ElMessage.success("删除成功");
          if (this.currentType && this.currentType.id === row.id) {
            this.currentType = null;
            this.dataList = [];
          }
          this.fetchTypes();
        });
      },
      handleDataCreate() {
        this.dataForm = this.getDefaultDataForm();
        this.dataForm.dictType = this.currentType.dictType;
        this.dataDialogVisible = true;
      },
      handleDataEdit(row) {
        this.dataForm = {
          id: row.id,
          dictType: row.dictType,
          label: row.label,
          value: row.value,
          tagClass: row.tagClass || "",
          sort: row.sort,
          status: row.status,
          remark: row.remark,
        };
        this.dataDialogVisible = true;
      },
      handleDataSubmit() {
        this.$refs.dataFormRef.validate(async (valid) => {
          if (!valid) return;
          this.submitting = true;
          try {
            if (this.dataForm.id) {
              await updateDictData(this.dataForm);
              ElMessage.success("修改成功");
            } else {
              await createDictData(this.dataForm);
              ElMessage.success("新增成功");
            }
            this.dataDialogVisible = false;
            this.fetchData();
          } finally {
            this.submitting = false;
          }
        });
      },
      handleDataDelete(row) {
        ElMessageBox.confirm(`确认删除字典数据 [${row.label}] 吗?`, "提示", {
          type: "warning",
        }).then(async () => {
          await deleteDictData({ ids: [row.id] });
          ElMessage.success("删除成功");
          this.fetchData();
        });
      },
    },
  };
</script>

<style lang="scss" scoped>
  .system-dict-container {
    padding: 15px;

    .table-toolbar {
      margin-bottom: 15px;
      display: flex;
      align-items: center;

      .current-type {
        margin-left: 12px;
        font-size: 13px;
        color: #909399;
      }
    }
  }
</style>
