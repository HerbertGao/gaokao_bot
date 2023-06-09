package com.herbertgao.telegram.bot.gaokao.database.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import java.time.LocalDateTime;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
@TableName(value = "gaokao.exam_date")
public class ExamDate implements Serializable {
    /**
     * ID
     */
    @TableId(value = "id", type = IdType.AUTO)
    private Integer id;

    /**
     * 考试年
     */
    @TableField(value = "exam_year")
    private Integer examYear;

    /**
     * 考试描述
     */
    @TableField(value = "exam_desc")
    private String examDesc;

    /**
     * 考试描述（短）
     */
    @TableField(value = "short_desc")
    private String shortDesc;

    /**
     * 考试开始时间
     */
    @TableField(value = "exam_begin_date")
    private LocalDateTime examBeginDate;

    /**
     * 考试结束时间
     */
    @TableField(value = "exam_end_date")
    private LocalDateTime examEndDate;

    /**
     * 考试年开始时间
     */
    @TableField(value = "exam_year_begin_date")
    private LocalDateTime examYearBeginDate;

    /**
     * 考试年结束时间
     */
    @TableField(value = "exam_year_end_date")
    private LocalDateTime examYearEndDate;

    /**
     * 是否删除
     */
    @TableField(value = "is_delete")
    private Boolean isDelete;

    private static final long serialVersionUID = 1L;

    public static final String COL_ID = "id";

    public static final String COL_EXAM_YEAR = "exam_year";

    public static final String COL_EXAM_DESC = "exam_desc";

    public static final String COL_SHORT_DESC = "short_desc";

    public static final String COL_EXAM_BEGIN_DATE = "exam_begin_date";

    public static final String COL_EXAM_END_DATE = "exam_end_date";

    public static final String COL_EXAM_YEAR_BEGIN_DATE = "exam_year_begin_date";

    public static final String COL_EXAM_YEAR_END_DATE = "exam_year_end_date";

    public static final String COL_IS_DELETE = "is_delete";
}