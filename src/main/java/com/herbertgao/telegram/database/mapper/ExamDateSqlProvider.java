package com.herbertgao.telegram.database.mapper;

import com.herbertgao.telegram.database.entity.ExamDate;
import com.herbertgao.telegram.database.entity.ExamDateExample.Criteria;
import com.herbertgao.telegram.database.entity.ExamDateExample.Criterion;
import com.herbertgao.telegram.database.entity.ExamDateExample;
import java.util.List;
import java.util.Map;
import org.apache.ibatis.jdbc.SQL;

public class ExamDateSqlProvider {

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String countByExample(ExamDateExample example) {
        SQL sql = new SQL();
        sql.SELECT("count(*)").FROM("exam_date");
        applyWhere(sql, example, false);
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String deleteByExample(ExamDateExample example) {
        SQL sql = new SQL();
        sql.DELETE_FROM("exam_date");
        applyWhere(sql, example, false);
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String insertSelective(ExamDate record) {
        SQL sql = new SQL();
        sql.INSERT_INTO("exam_date");
        
        if (record.getId() != null) {
            sql.VALUES("id", "#{id,jdbcType=INTEGER}");
        }
        
        if (record.getExamYear() != null) {
            sql.VALUES("exam_year", "#{examYear,jdbcType=INTEGER}");
        }
        
        if (record.getExamDesc() != null) {
            sql.VALUES("exam_desc", "#{examDesc,jdbcType=VARCHAR}");
        }
        
        if (record.getShortDesc() != null) {
            sql.VALUES("short_desc", "#{shortDesc,jdbcType=VARCHAR}");
        }
        
        if (record.getExamBeginDate() != null) {
            sql.VALUES("exam_begin_date", "#{examBeginDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamEndDate() != null) {
            sql.VALUES("exam_end_date", "#{examEndDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamYearBeginDate() != null) {
            sql.VALUES("exam_year_begin_date", "#{examYearBeginDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamYearEndDate() != null) {
            sql.VALUES("exam_year_end_date", "#{examYearEndDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getIsDelete() != null) {
            sql.VALUES("is_delete", "#{isDelete,jdbcType=BIT}");
        }
        
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String selectByExample(ExamDateExample example) {
        SQL sql = new SQL();
        if (example != null && example.isDistinct()) {
            sql.SELECT_DISTINCT("id");
        } else {
            sql.SELECT("id");
        }
        sql.SELECT("exam_year");
        sql.SELECT("exam_desc");
        sql.SELECT("short_desc");
        sql.SELECT("exam_begin_date");
        sql.SELECT("exam_end_date");
        sql.SELECT("exam_year_begin_date");
        sql.SELECT("exam_year_end_date");
        sql.SELECT("is_delete");
        sql.FROM("exam_date");
        applyWhere(sql, example, false);
        
        if (example != null && example.getOrderByClause() != null) {
            sql.ORDER_BY(example.getOrderByClause());
        }
        
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String updateByExampleSelective(Map<String, Object> parameter) {
        ExamDate record = (ExamDate) parameter.get("record");
        ExamDateExample example = (ExamDateExample) parameter.get("example");
        
        SQL sql = new SQL();
        sql.UPDATE("exam_date");
        
        if (record.getId() != null) {
            sql.SET("id = #{record.id,jdbcType=INTEGER}");
        }
        
        if (record.getExamYear() != null) {
            sql.SET("exam_year = #{record.examYear,jdbcType=INTEGER}");
        }
        
        if (record.getExamDesc() != null) {
            sql.SET("exam_desc = #{record.examDesc,jdbcType=VARCHAR}");
        }
        
        if (record.getShortDesc() != null) {
            sql.SET("short_desc = #{record.shortDesc,jdbcType=VARCHAR}");
        }
        
        if (record.getExamBeginDate() != null) {
            sql.SET("exam_begin_date = #{record.examBeginDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamEndDate() != null) {
            sql.SET("exam_end_date = #{record.examEndDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamYearBeginDate() != null) {
            sql.SET("exam_year_begin_date = #{record.examYearBeginDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamYearEndDate() != null) {
            sql.SET("exam_year_end_date = #{record.examYearEndDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getIsDelete() != null) {
            sql.SET("is_delete = #{record.isDelete,jdbcType=BIT}");
        }
        
        applyWhere(sql, example, true);
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String updateByExample(Map<String, Object> parameter) {
        SQL sql = new SQL();
        sql.UPDATE("exam_date");
        
        sql.SET("id = #{record.id,jdbcType=INTEGER}");
        sql.SET("exam_year = #{record.examYear,jdbcType=INTEGER}");
        sql.SET("exam_desc = #{record.examDesc,jdbcType=VARCHAR}");
        sql.SET("short_desc = #{record.shortDesc,jdbcType=VARCHAR}");
        sql.SET("exam_begin_date = #{record.examBeginDate,jdbcType=TIMESTAMP}");
        sql.SET("exam_end_date = #{record.examEndDate,jdbcType=TIMESTAMP}");
        sql.SET("exam_year_begin_date = #{record.examYearBeginDate,jdbcType=TIMESTAMP}");
        sql.SET("exam_year_end_date = #{record.examYearEndDate,jdbcType=TIMESTAMP}");
        sql.SET("is_delete = #{record.isDelete,jdbcType=BIT}");
        
        ExamDateExample example = (ExamDateExample) parameter.get("example");
        applyWhere(sql, example, true);
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    public String updateByPrimaryKeySelective(ExamDate record) {
        SQL sql = new SQL();
        sql.UPDATE("exam_date");
        
        if (record.getExamYear() != null) {
            sql.SET("exam_year = #{examYear,jdbcType=INTEGER}");
        }
        
        if (record.getExamDesc() != null) {
            sql.SET("exam_desc = #{examDesc,jdbcType=VARCHAR}");
        }
        
        if (record.getShortDesc() != null) {
            sql.SET("short_desc = #{shortDesc,jdbcType=VARCHAR}");
        }
        
        if (record.getExamBeginDate() != null) {
            sql.SET("exam_begin_date = #{examBeginDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamEndDate() != null) {
            sql.SET("exam_end_date = #{examEndDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamYearBeginDate() != null) {
            sql.SET("exam_year_begin_date = #{examYearBeginDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getExamYearEndDate() != null) {
            sql.SET("exam_year_end_date = #{examYearEndDate,jdbcType=TIMESTAMP}");
        }
        
        if (record.getIsDelete() != null) {
            sql.SET("is_delete = #{isDelete,jdbcType=BIT}");
        }
        
        sql.WHERE("id = #{id,jdbcType=INTEGER}");
        
        return sql.toString();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method corresponds to the database table exam_date
     *
     * @mbg.generated Thu Sep 05 13:16:58 CST 2019
     */
    protected void applyWhere(SQL sql, ExamDateExample example, boolean includeExamplePhrase) {
        if (example == null) {
            return;
        }
        
        String parmPhrase1;
        String parmPhrase1_th;
        String parmPhrase2;
        String parmPhrase2_th;
        String parmPhrase3;
        String parmPhrase3_th;
        if (includeExamplePhrase) {
            parmPhrase1 = "%s #{example.oredCriteria[%d].allCriteria[%d].value}";
            parmPhrase1_th = "%s #{example.oredCriteria[%d].allCriteria[%d].value,typeHandler=%s}";
            parmPhrase2 = "%s #{example.oredCriteria[%d].allCriteria[%d].value} and #{example.oredCriteria[%d].criteria[%d].secondValue}";
            parmPhrase2_th = "%s #{example.oredCriteria[%d].allCriteria[%d].value,typeHandler=%s} and #{example.oredCriteria[%d].criteria[%d].secondValue,typeHandler=%s}";
            parmPhrase3 = "#{example.oredCriteria[%d].allCriteria[%d].value[%d]}";
            parmPhrase3_th = "#{example.oredCriteria[%d].allCriteria[%d].value[%d],typeHandler=%s}";
        } else {
            parmPhrase1 = "%s #{oredCriteria[%d].allCriteria[%d].value}";
            parmPhrase1_th = "%s #{oredCriteria[%d].allCriteria[%d].value,typeHandler=%s}";
            parmPhrase2 = "%s #{oredCriteria[%d].allCriteria[%d].value} and #{oredCriteria[%d].criteria[%d].secondValue}";
            parmPhrase2_th = "%s #{oredCriteria[%d].allCriteria[%d].value,typeHandler=%s} and #{oredCriteria[%d].criteria[%d].secondValue,typeHandler=%s}";
            parmPhrase3 = "#{oredCriteria[%d].allCriteria[%d].value[%d]}";
            parmPhrase3_th = "#{oredCriteria[%d].allCriteria[%d].value[%d],typeHandler=%s}";
        }
        
        StringBuilder sb = new StringBuilder();
        List<Criteria> oredCriteria = example.getOredCriteria();
        boolean firstCriteria = true;
        for (int i = 0; i < oredCriteria.size(); i++) {
            Criteria criteria = oredCriteria.get(i);
            if (criteria.isValid()) {
                if (firstCriteria) {
                    firstCriteria = false;
                } else {
                    sb.append(" or ");
                }
                
                sb.append('(');
                List<Criterion> criterions = criteria.getAllCriteria();
                boolean firstCriterion = true;
                for (int j = 0; j < criterions.size(); j++) {
                    Criterion criterion = criterions.get(j);
                    if (firstCriterion) {
                        firstCriterion = false;
                    } else {
                        sb.append(" and ");
                    }
                    
                    if (criterion.isNoValue()) {
                        sb.append(criterion.getCondition());
                    } else if (criterion.isSingleValue()) {
                        if (criterion.getTypeHandler() == null) {
                            sb.append(String.format(parmPhrase1, criterion.getCondition(), i, j));
                        } else {
                            sb.append(String.format(parmPhrase1_th, criterion.getCondition(), i, j,criterion.getTypeHandler()));
                        }
                    } else if (criterion.isBetweenValue()) {
                        if (criterion.getTypeHandler() == null) {
                            sb.append(String.format(parmPhrase2, criterion.getCondition(), i, j, i, j));
                        } else {
                            sb.append(String.format(parmPhrase2_th, criterion.getCondition(), i, j, criterion.getTypeHandler(), i, j, criterion.getTypeHandler()));
                        }
                    } else if (criterion.isListValue()) {
                        sb.append(criterion.getCondition());
                        sb.append(" (");
                        List<?> listItems = (List<?>) criterion.getValue();
                        boolean comma = false;
                        for (int k = 0; k < listItems.size(); k++) {
                            if (comma) {
                                sb.append(", ");
                            } else {
                                comma = true;
                            }
                            if (criterion.getTypeHandler() == null) {
                                sb.append(String.format(parmPhrase3, i, j, k));
                            } else {
                                sb.append(String.format(parmPhrase3_th, i, j, k, criterion.getTypeHandler()));
                            }
                        }
                        sb.append(')');
                    }
                }
                sb.append(')');
            }
        }
        
        if (sb.length() > 0) {
            sql.WHERE(sb.toString());
        }
    }
}